# Decimal wire contract

Amounts in accounting.v1 are decimal strings, never JSON/Protobuf floats.
For monetary amounts the accepted domain scale is two decimal places, bounded
by the Ledger Money range (absolute minor units <= 2^50). Posting amounts must
be positive. No implicit monetary rounding is allowed at the transport boundary.
Canonical amounts use a dot decimal separator and exactly two fractional digits,
with no exponent, grouping separators, whitespace, leading plus or negative zero.

Quantity and tax rate are decimal inputs to invoice calculation, not Money.
Their range/scale and invoice tax rounding remain governed by the invoice domain;
they must not be silently subjected to a two-decimal monetary conversion.
The current Ledger compatibility adapter still converts these strings to float64.
A lossless parser and a versioned domain migration remain necessary before claiming
end-to-end decimal precision. Existing wire fields and numbers must be preserved.

iam.v1.IdentityService defines the authentication, authorization, role, and licensing boundary. Profiles and contact lifecycle belong to user.v1.ProfileService.

## Ledger RPC input precision (2026-10-12)

Ledger monetary input is plain decimal text at scale 10^-2. Invoice quantity, unit price and fractional tax rate use scale 10^-4 (23% is `0.2300`). Nonzero excess precision, exponent notation and whitespace are rejected; optional signs and extra trailing fractional zeroes are accepted. Quantity must be positive, price nonnegative, tax rate between zero and one. Scaled absolute input values cannot exceed 2^50; invoice totals additionally cannot exceed 2^50 minor units. This is the implemented Ledger policy, not a statement that every invoice RPC in another service is implemented.
