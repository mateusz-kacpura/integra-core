/**
 * ============================================================================
 * Integra Platform Docs — Interactive Pan, Zoom & Fullscreen for Diagrams
 * Features:
 *   - Mouse Drag-to-Pan & Mouse Wheel Zoom (centered on cursor)
 *   - Touch 1-finger Pan & 2-finger Pinch-to-Zoom
 *   - Floating Glassmorphism Toolbar with Zoom In, Zoom Out, Reset & Fullscreen
 *   - Real-time Zoom Percentage Indicator
 *   - Double-Click to toggle Fullscreen, ESC to exit
 *   - Full compatibility with Material for MkDocs (Shadow DOM, Slate/Light, Instant Loading)
 * ============================================================================
 */

(function () {
  'use strict';

  // 1. Intercept attachShadow to ensure Mermaid's Shadow DOM is open and accessible
  const origAttachShadow = Element.prototype.attachShadow;
  Element.prototype.attachShadow = function (options) {
    const opts = Object.assign({}, options, { mode: 'open' });
    const shadow = origAttachShadow.call(this, opts);
    this.__diagramShadowRoot = shadow;
    return shadow;
  };

  const ICONS = {
    zoomIn: '<svg viewBox="0 0 24 24"><path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/></svg>',
    zoomOut: '<svg viewBox="0 0 24 24"><path d="M19 13H5v-2h14v2z"/></svg>',
    reset: '<svg viewBox="0 0 24 24"><path d="M12 5V1L7 6l5 5V7c3.31 0 6 2.69 6 6s-2.69 6-6 6-6-2.69-6-6H4c0 4.42 3.58 8 8 8s8-3.58 8-8-3.58-8-8-8z"/></svg>',
    fullscreen: '<svg viewBox="0 0 24 24"><path d="M7 14H5v5h5v-2H7v-3zm-2-4h2V7h3V5H5v5zm12 7h-3v2h5v-5h-2v3zM14 5v2h3v3h2V5h-5z"/></svg>',
    exitFullscreen: '<svg viewBox="0 0 24 24"><path d="M5 16h3v3h2v-5H5v2zm3-8H5v2h5V5H8v3zm6 11h2v-3h3v-2h-5v5zm2-14v3h3v2h-5V5h2z"/></svg>'
  };

  /**
   * Initializes a diagram element with interactive zoom, pan, and fullscreen.
   */
  function initInteractiveDiagram(diagramEl) {
    // Prevent duplicate initialization
    if (diagramEl.dataset.diagramZoomDone || diagramEl.closest('.diagram-zoom-wrapper')) {
      return;
    }

    // Must be a rendered diagram (either has shadowRoot with SVG or contains SVG, or is div.mermaid)
    const shadow = diagramEl.shadowRoot || diagramEl.__diagramShadowRoot;
    const svg = shadow ? shadow.querySelector('svg') : diagramEl.querySelector('svg');

    // If it's a pre.mermaid that hasn't rendered yet, let the observer handle it when replaced
    if (!svg && diagramEl.tagName.toLowerCase() === 'pre') {
      return;
    }

    diagramEl.dataset.diagramZoomDone = 'true';

    // Ensure internal SVG styles allow unconstrained zooming
    if (svg) {
      svg.style.maxWidth = 'none';
      svg.style.overflow = 'visible';
    }

    // Build DOM structure: wrapper > (toolbar + viewport > stage > diagramEl) + hint
    const wrapper = document.createElement('div');
    wrapper.className = 'diagram-zoom-wrapper';

    const viewport = document.createElement('div');
    viewport.className = 'diagram-zoom-viewport';
    viewport.setAttribute('tabindex', '0');
    viewport.setAttribute('role', 'region');
    viewport.setAttribute('aria-label', 'Interaktywny diagram');

    const stage = document.createElement('div');
    stage.className = 'diagram-zoom-stage';

    // Insert wrapper into DOM where diagramEl was
    diagramEl.parentNode.insertBefore(wrapper, diagramEl);
    stage.appendChild(diagramEl);
    viewport.appendChild(stage);
    wrapper.appendChild(viewport);

    // State
    let scale = 1.0;
    let translateX = 0;
    let translateY = 0;
    let isDragging = false;
    let startX = 0;
    let startY = 0;

    const MIN_SCALE = 0.2;
    const MAX_SCALE = 12.0;

    // Scale label
    const scaleLabel = document.createElement('span');
    scaleLabel.className = 'diagram-zoom-scale';
    scaleLabel.textContent = '100%';

    function updateTransform() {
      stage.style.transform = `translate(${translateX}px, ${translateY}px) scale(${scale})`;
      scaleLabel.textContent = `${Math.round(scale * 100)}%`;
    }

    function setZoom(newScale, focalX, focalY) {
      const clamped = Math.min(Math.max(newScale, MIN_SCALE), MAX_SCALE);
      if (clamped === scale) return;

      if (focalX !== undefined && focalY !== undefined) {
        const rect = viewport.getBoundingClientRect();
        const offsetX = focalX - rect.left - rect.width / 2;
        const offsetY = focalY - rect.top - rect.height / 2;
        translateX -= (offsetX / scale) * (clamped - scale);
        translateY -= (offsetY / scale) * (clamped - scale);
      }

      scale = clamped;
      updateTransform();
    }

    function resetView() {
      scale = 1.0;
      translateX = 0;
      translateY = 0;
      updateTransform();
    }

    // Build Floating Toolbar
    const toolbar = document.createElement('div');
    toolbar.className = 'diagram-zoom-toolbar';
    toolbar.setAttribute('role', 'toolbar');
    toolbar.setAttribute('aria-label', 'Narzędzia powiększenia i pełnego ekranu');

    const btnZoomIn = createBtn(ICONS.zoomIn, 'Przybliż (+)', () => setZoom(scale * 1.25));
    const btnZoomOut = createBtn(ICONS.zoomOut, 'Oddal (-)', () => setZoom(scale * 0.8));
    const btnReset = createBtn(ICONS.reset, 'Zresetuj widok i wyśrodkuj', resetView);
    const btnFullscreen = createBtn(ICONS.fullscreen, 'Tryb pełnoekranowy (Podwójny klik lub ESC)', toggleFullscreen);

    toolbar.appendChild(btnZoomIn);
    toolbar.appendChild(scaleLabel);
    toolbar.appendChild(btnZoomOut);
    toolbar.appendChild(btnReset);
    toolbar.appendChild(btnFullscreen);
    wrapper.appendChild(toolbar);

    // Floating Hint
    const hint = document.createElement('div');
    hint.className = 'diagram-zoom-hint';
    hint.innerHTML = '<span>Kółko: <b>zoom</b> &bull; Przeciągnij: <b>przesuń</b> &bull; Podwójny klik: <b>pełny ekran</b></span>';
    wrapper.appendChild(hint);

    // Fullscreen logic
    function toggleFullscreen() {
      const isFull = wrapper.classList.toggle('is-fullscreen');
      btnFullscreen.innerHTML = isFull ? ICONS.exitFullscreen : ICONS.fullscreen;
      btnFullscreen.title = isFull ? 'Zamknij pełny ekran (ESC)' : 'Tryb pełnoekranowy';
      resetView();

      if (isFull) {
        document.body.style.overflow = 'hidden';
      } else {
        document.body.style.overflow = '';
      }
    }

    // ESC to exit fullscreen
    window.addEventListener('keydown', (e) => {
      if (e.key === 'Escape' && wrapper.classList.contains('is-fullscreen')) {
        toggleFullscreen();
      }
    });

    // Mouse drag to pan
    viewport.addEventListener('mousedown', (e) => {
      if (e.target.closest('.diagram-zoom-toolbar')) return;
      isDragging = true;
      viewport.classList.add('is-dragging');
      startX = e.clientX - translateX;
      startY = e.clientY - translateY;
      e.preventDefault();
    });

    window.addEventListener('mousemove', (e) => {
      if (!isDragging) return;
      translateX = e.clientX - startX;
      translateY = e.clientY - startY;
      updateTransform();
    });

    window.addEventListener('mouseup', () => {
      if (isDragging) {
        isDragging = false;
        viewport.classList.remove('is-dragging');
      }
    });

    // Mouse wheel zoom (centered at cursor)
    viewport.addEventListener('wheel', (e) => {
      e.preventDefault();
      const zoomFactor = e.deltaY < 0 ? 1.15 : 0.87;
      setZoom(scale * zoomFactor, e.clientX, e.clientY);
    }, { passive: false });

    // Double-click toggle fullscreen
    viewport.addEventListener('dblclick', (e) => {
      if (e.target.closest('.diagram-zoom-toolbar')) return;
      toggleFullscreen();
    });

    // Touch support (1-finger pan, 2-finger pinch)
    let touchStartX = 0;
    let touchStartY = 0;
    let initialPinchDist = null;

    viewport.addEventListener('touchstart', (e) => {
      if (e.touches.length === 1) {
        isDragging = true;
        touchStartX = e.touches[0].clientX - translateX;
        touchStartY = e.touches[0].clientY - translateY;
      } else if (e.touches.length === 2) {
        isDragging = false;
        initialPinchDist = Math.hypot(
          e.touches[0].clientX - e.touches[1].clientX,
          e.touches[0].clientY - e.touches[1].clientY
        );
      }
    }, { passive: true });

    viewport.addEventListener('touchmove', (e) => {
      if (isDragging && e.touches.length === 1) {
        translateX = e.touches[0].clientX - touchStartX;
        translateY = e.touches[0].clientY - touchStartY;
        updateTransform();
      } else if (e.touches.length === 2 && initialPinchDist) {
        const currentDist = Math.hypot(
          e.touches[0].clientX - e.touches[1].clientX,
          e.touches[0].clientY - e.touches[1].clientY
        );
        const factor = currentDist / initialPinchDist;
        if (factor > 1.03 || factor < 0.97) {
          const midX = (e.touches[0].clientX + e.touches[1].clientX) / 2;
          const midY = (e.touches[0].clientY + e.touches[1].clientY) / 2;
          setZoom(scale * factor, midX, midY);
          initialPinchDist = currentDist;
        }
      }
    }, { passive: true });

    viewport.addEventListener('touchend', () => {
      isDragging = false;
      initialPinchDist = null;
    });
  }

  function createBtn(iconHtml, tooltip, onClick) {
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'diagram-zoom-btn';
    btn.title = tooltip;
    btn.setAttribute('aria-label', tooltip);
    btn.innerHTML = iconHtml;
    btn.addEventListener('click', (e) => {
      e.stopPropagation();
      onClick();
    });
    return btn;
  }

  /**
   * Scans document for rendered Mermaid blocks or diagram SVGs
   */
  function scanAndInit() {
    // 1. Check for div.mermaid (created by Material after rendering Mermaid)
    const divMermaids = document.querySelectorAll('div.mermaid');
    divMermaids.forEach((el) => {
      initInteractiveDiagram(el);
    });

    // 2. Also check any pre.mermaid that might already contain an SVG
    const preMermaids = document.querySelectorAll('pre.mermaid');
    preMermaids.forEach((el) => {
      const svg = el.querySelector('svg');
      if (svg) {
        initInteractiveDiagram(el);
      }
    });
  }

  // Set up MutationObserver to detect when Material replaces pre.mermaid with div.mermaid
  function observeContent() {
    const content = document.querySelector('.md-content') || document.body;
    if (!content) return;

    const observer = new MutationObserver((mutations) => {
      let needsScan = false;
      for (const m of mutations) {
        if (m.addedNodes.length > 0) {
          for (const node of m.addedNodes) {
            if (node.nodeType === Node.ELEMENT_NODE) {
              if (
                node.classList && (node.classList.contains('mermaid') || node.querySelector('.mermaid'))
              ) {
                needsScan = true;
                break;
              }
            }
          }
        }
        if (needsScan) break;
      }
      if (needsScan) {
        setTimeout(scanAndInit, 50);
      }
    });

    observer.observe(content, { childList: true, subtree: true });
  }

  // Register hooks with Material for MkDocs instant loading
  function register() {
    scanAndInit();
    observeContent();

    // Repeat scan at intervals to catch late async Mermaid CDN renders
    setTimeout(scanAndInit, 200);
    setTimeout(scanAndInit, 600);
    setTimeout(scanAndInit, 1200);
    setTimeout(scanAndInit, 2500);

    // MkDocs Material navigation observable (instant loading)
    if (typeof window.document$ !== 'undefined') {
      window.document$.subscribe(() => {
        setTimeout(scanAndInit, 100);
        setTimeout(scanAndInit, 500);
        setTimeout(scanAndInit, 1200);
      });
    }
  }

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', register);
  } else {
    register();
  }

  window.addEventListener('load', () => {
    setTimeout(scanAndInit, 300);
    setTimeout(scanAndInit, 1000);
  });
})();
