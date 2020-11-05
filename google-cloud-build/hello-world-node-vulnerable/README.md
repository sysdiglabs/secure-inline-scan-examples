# Example Node.js hello world containerized application

This repository contains the official Node.js "hello world" containerized application files, as [described in this link on their website](https://nodejs.org/de/docs/guides/nodejs-docker-webapp/).

It is an useful example to look for vulnerabilities using [Sysdig Secure image scanning](https://www.sysdig.com).

To make the image secure, change the reference base image from `node:12` to `bitnami/node`.
