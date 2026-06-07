# multi_root

Demonstrates multiple source roots (YZC-0022).

`proj/` is the project directory (owns `target/`). `lib/` is a separate
source root contributing the `greeting` boc to the same FQN namespace.

```
lib/greeting.yz   →  FQN: greeting
proj/main.yz      →  FQN: main
```

## Run

```
yzc build proj/ lib/
proj/target/bin/app
```

## Output

```
Hello from lib!
Goodbye from lib!
```
