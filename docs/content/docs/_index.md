---
title: "Documentation"
weight: 1
---

Welcome to the documentation for **`pd`** (go-pixeldrain), a fast and versatile command-line interface, Model Context Protocol (MCP) server, and Go SDK for [Pixeldrain](https://pixeldrain.com/).

## Overview

Pixeldrain is a reliable and fast file-sharing and cloud storage service. `pd` brings the full power of Pixeldrain into your terminal, scripts, and AI workflows.

### Key Capabilities

* **Command Line Interface (`pd`)**: Upload, download, create albums, and manage files without touching a web browser.
* **Model Context Protocol (`pd mcp`)**: Seamlessly connect Pixeldrain file capabilities to AI agents (Claude Desktop, Cursor, etc.).
* **End-to-End Encryption**: Encrypt files client-side before sending them to the cloud using [age](https://github.com/FiloSottile/age) or SSH keys.
* **Unix Pipelines**: Native support for `stdin`, `stdout`, and `tar` streaming.
* **Go SDK**: Programmatic access to the Pixeldrain REST API v2 in Go applications.

## Table of Contents

* [Quick Start](quickstart): Learn how to install `pd` and start sharing files immediately.
* [MCP Server (AI Integration)](mcp): Configure AI agents to transfer files through Pixeldrain.
* [End-to-End Encryption](encryption): Encrypt your uploads client-side for maximum confidentiality.
* [Go SDK](sdk): Use `go-pixeldrain` as a library in your Go projects.
