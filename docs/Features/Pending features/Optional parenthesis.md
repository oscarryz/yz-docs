---
type: pending-feature
generated: { by: "oscarryz", at: 2023-12-06T20:13:02-06:00 }
---
#pending-feature

When there's at least one arg and no named params

    a.sum b

    a.print 'A'

vs

    a.sum(b)
    a.print('A')


No args? Then still needs the params

Style choice probably enforced by tool, use parenthesis always.