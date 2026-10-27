---
title: "pd (Pixeldrain CLI & MCP Server)"
layout: "hextra-home"
toc: false
---

<div class="hx-text-center hx-w-full">

<div class="hx-mt-6 hx-mb-6">
{{< hextra/hero-headline >}}
PixelDrain CLI & MCP Server
{{< /hextra/hero-headline >}}
</div>

<div class="hx-mb-6">
{{< hextra/hero-badge link="https://github.com/jkawamoto/go-pixeldrain/releases" >}}
  <div class="hx-w-2 hx-h-2 hx-rounded-full hx-bg-primary-400"></div>
  <span>Latest Release v0.9.0 &middot; MCP Ready</span>
  {{< icon name="arrow-circle-right" attributes="height=14" >}}
{{< /hextra/hero-badge >}}
</div>

<div class="hx-mb-16">
{{< hextra/hero-subtitle >}}
Upload, download, and share files, media, and albums directly from your terminal,&nbsp;<br class="sm:block hidden" />shell scripts, pipelines, and AI workflows without opening a browser.
{{< /hextra/hero-subtitle >}}
</div>

<div class="hx-mb-16">
  {{< hextra/hero-button text="Get Started" link="docs" >}}
</div>

</div>

<div class="hx-mb-12">
  <img src="images/pd-banner.jpg" alt="pd - Pixeldrain CLI & MCP tool" class="hx-w-full hx-rounded-xl hx-shadow-lg" />
</div>

{{< hextra/feature-grid >}}
  {{< hextra/feature-card
    title="Instant File Sharing"
    subtitle="Upload files, screenshots, or generated images from the CLI and get a shareable Pixeldrain link instantly."
    icon="lightning-bolt"
    link="docs/quickstart"
  >}}

  {{< hextra/feature-card
    title="AI & MCP Integration"
    subtitle="Built-in Model Context Protocol server (<code>pd mcp</code>) enabling Claude Desktop, Cursor, and agents to transfer files directly."
    icon="chip"
    link="docs/mcp"
  >}}

  {{< hextra/feature-card
    title="Album & Pipeline Friendly"
    subtitle="Bundle multiple files into albums with <code>--album</code>, pipe from <code>stdin</code>, or stream directory archives with <code>tar</code>."
    icon="folder"
    link="docs/quickstart"
  >}}

  {{< hextra/feature-card
    title="End-to-End Encryption"
    subtitle="Client-side file encryption with age or SSH public keys before uploading, ensuring maximum privacy."
    icon="shield-check"
    link="docs/encryption"
  >}}
{{< /hextra/feature-grid >}}

<div class="hx-mt-16"></div>

## Quick Installation

### Homebrew (macOS & Linux)

```shell
brew tap jkawamoto/pixeldrain
brew install pixeldrain
```

### Go Install

```shell
go install github.com/jkawamoto/go-pixeldrain/cmd/pd@latest
```

Pre-compiled binaries for macOS, Linux, and Windows are also available on [GitHub Releases](https://github.com/jkawamoto/go-pixeldrain/releases).

<div class="hx-mt-12"></div>

## Quick Start in 30 Seconds

```shell
# Upload a file and get a shareable link
pd upload artwork.png
# => https://pixeldrain.com/u/xxxxxx

# Upload multiple files as a single album
pd upload --album "My Project" screenshot1.png screenshot2.png
# => https://pixeldrain.com/l/xxxxxx

# Download a file or full album
pd download https://pixeldrain.com/u/xxxxxx
pd download --all https://pixeldrain.com/l/xxxxxx

# Start MCP server for AI tools
pd mcp
```
