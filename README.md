<p align="center"><img src="docs/images/logo_v1.png" width="250"></p>

<p align="center"><a href="https://github.com/afritzler/oopsie/actions"><img src="https://github.com/afritzler/oopsie/workflows/Docker/badge.svg"></a></p>


# oopsie

Oopsie [/ˈuːpsi/] is a Kubernetes controller that watches all `Events` within a cluster and enriches failed objects with solutions found on [stackoverflow](https://stackoverflow.com).

## Why Oopsie?

Kubernetes is a great tool for orchestrating containerized workloads on a fleet of machines. Unfortunatelly, it is sometimes not that easy for new Kubernetes users to resolve problems which occur in their deployments. The illustration below is a visual representation of what happens when you deploy your first application after having mastered the [Wordpress and Guestbook example](https://github.com/kubernetes/examples).

<p align="center"><img src="docs/images/shipit.gif" width="200"></p>

## Installation