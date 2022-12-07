[![Go](https://github.com/retr0h/task-migrate/actions/workflows/go.yml/badge.svg)](https://github.com/retr0h/task-migrate/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/retr0h/task-migrate/branch/main/graph/badge.svg?token=9OYBYGKKJD)](https://codecov.io/gh/retr0h/task-migrate)

# Task Migrate

A general purpose migration tool built around [Task][] runner,
where migration files are simply [Task][] targets.

Define "sql migration" like [Task][] files in a `versions` directory. These
[Task][] files will be executed in lexicographical order.  Tasks which
were successfully executed will not be executed again.

## Purpose

Intended to aid in the upgrading of embeded systems Kubernetes' versions.

## Usage

### Configuration

Create a migration file `1_delete_prom_for_upgrade.yaml`:

```yaml
---
version: "3"

vars:
  K3S_VERSION:
    sh: k3d version | grep k3s | awk '{print $3}'
  TO_K3S_VERSION: v1.23.1-k3s1

tasks:
  up:
    cmds:
      - kubectl delete deployment kube-prometheus-stack-kube-state-metrics -n monitoring
      - kubectl delete ds kube-prometheus-stack-prometheus-node-exporter -n monitoring
    preconditions:
      - sh: "[ {{ .K3S_VERSION }} == {{ .TO_K3S_VERSION }} ]"
        msg: Kubernetes to version not matched
```

### Status

Display migration status:

```bash
$ go run main.go status
+-----------+---------+
| MIGRATION | APPLIED |
+-----------+---------+
+-----------+---------+
```

### Up

Execute all migrations:

```bash
$ go run main.go -d testdata/versions up

task: [up] echo "1 - UP"
1 - UP
task: [up] echo "2 - UP"
2 - UP
```

Dislay migration status:

```bash
$ go run main.go status
+-----------+----------------------------------------+
| MIGRATION | APPLIED                                |
+-----------+----------------------------------------+
| 1-foo.yml | 2022-12-06 15:00:47.493033 -0800 -0800 |
| 2-bar.yml | 2022-12-06 15:00:47.494593 -0800 -0800 |
+-----------+----------------------------------------+
```

## Testing

Run unit tests:

```sh
task go:test
```

## License

The [MIT] License.

[MIT]: LICENSE
[Task]: https://github.com/go-task/task
