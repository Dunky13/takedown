# Google terms and takedown: research findings

> Historical report: this describes the original browser-automation implementation before the documented-API redesign. Current behavior and acceptance are recorded in [AUTOMATION.md](AUTOMATION.md). Some source links below refer to files removed by that redesign and are no longer available in the checkout.

Accessed: **2026-10-25** for every external source linked below. Scope: independent public-source research and local source inspection only. No account login, account data access, app execution, or app source modification. This report distinguishes contractual language, support policy, and inference; it does not determine enforceability under national law.

## Practical verdict

**Do not describe the present implementation as clearly permitted by Google.** Downloading the user's own Takeout data is supported. Automated Google sign-in conflicts with Google's supported-browser guidance. Stealth creates a credible **conditional anti-bypass risk**: if used to defeat Google's protective measures, the terms' anti-bypass restriction applies. Mere automation, DOM inspection, or possession of one's own data does not, by itself, establish a contractual breach or exemption. Sources and limits follow.

## What the app actually does

The current checkout confirms the supplied description:

| Behavior | Source evidence |
| --- | --- |
| Installed Chromium under Rod control, headless operations, restored cookies/profile, and `stealth.Page` | [browser.go](../internal/auth/browser.go), lines 44-111; [go.mod](../go.mod) pins stealth v0.4.9. |
| Automatic email/password and optional TOTP submission; visible fallback for other challenges | [login.go](../internal/auth/login.go), lines 33-95; [browser.go](../internal/auth/browser.go), lines 228-355. The visible browser remains controlled. |
| Parse existing exports and capture signed download URLs through CDP request/response interception | [parser.go](../internal/auth/parser.go); [browser.go](../internal/auth/browser.go), lines 503-527; [intercept.go](../internal/auth/intercept.go), lines 21-76. |
| Separate HTTP downloads, including a possible one-byte metadata probe | [http.go](../internal/auth/http.go), lines 188-205; [manager.go](../internal/auth/manager.go), lines 130-175. No OAuth integration was found in the inspected application paths. |

The dependency's [immutable v0.4.9 source](https://raw.githubusercontent.com/go-rod/stealth/v0.4.9/main.go) explicitly describes anti-bot detection and injects JavaScript on new documents. That establishes a concealment mechanism, not proof of which Google protection it defeats.

## Contractual terms, including Netherlands/EEA

The [Netherlands country version](https://policies.google.com/terms?hl=en&gl=NL) is effective **July 30, 2026**, with Google Ireland Limited as provider for the EEA and Switzerland. The [US version](https://policies.google.com/terms?hl=en&gl=us) has the same effective date. The May 22, 2024 text is archived.

Relevant current provisions prohibit bypassing Google's systems or protective measures, deceptive access, and automated access that violates machine-readable instructions. The automation clause is conditional, not a universal prohibition on every automated request. The identity-concealment provision also includes a purpose condition: concealing identity to violate the terms. Content ownership does not remove access-method conditions. The Netherlands text preserves mandatory legal rights and generally applies the resident's national law. [Netherlands terms](https://policies.google.com/terms?hl=en&gl=NL).

**Application, by inference:** stealth-assisted automated authentication could constitute prohibited bypassing if it circumvents an actual protection. Neither successful authentication nor a downloader's benign purpose resolves that question. CDP interception is not expressly named as categorically forbidden in these terms. No actual protective-measure circumvention was verified.

## Sign-in policy versus OAuth-specific rules

Google's current [supported-browser help](https://support.google.com/accounts/answer/7675428?hl=en&co=GENIE.Platform%3DDesktop) says Google might stop sign-ins from browsers controlled through software automation or embedded in another application. This directly concerns the app's authentication method, but is a support/security policy, not an independently demonstrated contractual breach. Headed Chrome and manual phone/passkey completion do not remove Rod control.

Google's [August 28, 2020 developer guidance](https://developers.googleblog.com/en/guidance-to-developers-affected-by-our-effort-to-block-less-secure-browsers-and-applications/) describes browser requirements against headless operation, automated sign-in, request replay/network alteration, and misleading browser identification. It announced embedded-framework blocking from January 4, 2021. This is historical guidance, largely explaining supported OAuth/browser integration; current account help independently confirms automation-related blocking.

The [OAuth 2.0 Policies](https://developers.google.com/identity/protocols/oauth2/policies), last modified **August 5, 2026**, apply to developers using OAuth, including OpenID Connect. Their secure-browser rule prohibits directing Google OAuth authorization requests into developer-controlled embedded agents, including agents able to inject scripts or access cookies. **No OAuth client/authorization integration was found here, so that rule is not established as a direct violation by this app.** Ordinary account-browser restrictions remain relevant. An API redesign would bring OAuth rules directly into scope.

## Robots: actual directives and limited implications

The live [Takeout robots file](https://takeout.google.com/robots.txt) returned HTTP 200:

```text
User-agent: *
Allow: *
Disallow: */_/*
Disallow: /u/*
```

The fixed management route `/settings/takeout/downloads` and accepted download routes `/takeout/download` and `/settings/takeout/download` do not themselves match those disallows. Runtime redirects, account-prefixed routes, and subrequests might. **No live runtime URL trace or robots-enforcement audit was performed.** [manager.go](../internal/auth/manager.go), line 21; [parser.go](../internal/auth/parser.go), line 145.

[Accounts robots](https://accounts.google.com/robots.txt) allows the root while excluding several named legacy/authentication paths; it does not categorically disallow every modern sign-in route. The Takeout download host's [robots URL](https://takeout.googleusercontent.com/robots.txt) returned HTTP 404. Allowance or absence is not blanket permission and does not override protective-measure restrictions or sign-in policy.

## Takeout's operational limits

[Takeout help](https://support.google.com/accounts/answer/3024190?hl=en) explicitly supports exporting and downloading one's data. Each archive expires in approximately **seven days** and may be downloaded **five times**; Google directs users to request another archive afterward. Re-entering a password and completing additional verification are documented protections.

**Unresolved counting:** the help does not define whether range probes, resumed requests, redirects, interrupted downloads, or repeated link retrieval consume that allowance. Do not promise that the app's probe is free, that each successful completed file is the counting unit, or that refreshing a URL resets an archive's allowance. Signed-link lifetime and archive lifetime should not be conflated.

## Official alternatives and their tradeoffs

1. **Recommended: human-operated Takeout plus a download/import workflow.** Official delivery options include a browser download and export to Drive, Dropbox, OneDrive, or Box. The help also documents scheduled exports and URL parameters that preselect products, destination, and frequency. These documented integration points do not authorize automated credential entry or arbitrary scraping. They can require a new export if existing archives were created with another destination. [Takeout help](https://support.google.com/accounts/answer/3024190?hl=en).
2. **Drive-delivered archives plus an authorized Drive API downloader.** Google documents binary downloads with `files.get` and `alt=media`, with read-authorizing OAuth scopes. This changes the delivery/authentication design and requires permission to the selected files. [Drive download guide](https://developers.google.com/workspace/drive/api/guides/manage-downloads).
3. **Data Portability API for supported new exports.** The [API reference](https://developers.google.com/data-portability/reference/rest) exposes archive initiation and job lookup, not a documented list/download interface for pre-existing web Takeout exports. Treat that as a documented-interface finding, not proof that no private endpoint exists. It is **not a drop-in replacement** for the present requirement.

Data Portability constraints:

- Availability is regional, including the Netherlands; managed accounts, under-18 accounts, and Advanced Protection accounts may be ineligible. [Account help](https://support.google.com/accounts/answer/14452558).
- Supported products are limited to the published [OAuth scope list](https://developers.google.com/data-portability/user-guide/scopes), not arbitrary Takeout content. Public release normally requires verification; restricted scopes can require security assessment. [Overview](https://developers.google.com/data-portability/user-guide/overview).
- Grants can be one-time or 30/180 days; one-time means one export per scope, while time-based access permits exports every 24 hours. [Workflow](https://developers.google.com/data-portability/user-guide/introduction).
- API download URLs expire after six hours; exported data remains available for 14 days. These are separate from web Takeout's limits. [Methods](https://developers.google.com/data-portability/user-guide/methods).
- Its [developer policy](https://developers.google.com/data-portability/policy) requires a user-benefiting transfer use case, minimum permissions, transparent disclosures, and secure data handling. Approval for this app's exact design is unverified.

## Recommendation and remaining uncertainty

Prefer an official delivery/API workflow and remove automated Google credential/TOTP submission and stealth from a redesigned integration. Simply removing stealth would reduce concealment risk without establishing permission for the remaining browser automation. For continued automated access to existing web archives, seek explicit Google clarification before advertising compliance.

Unresolved: actual account country/terms, account-specific restrictions, runtime routes and protections, download counting, and approval eligibility. No live-account acceptance or enforcement outcome was tested.

Retrieval limits: the web research tool refused Netherlands query URLs and Takeout robots; unauthenticated direct public HTTP retrieval recovered both, including the country label and effective date. The download-host robots URL was genuinely 404. No authenticated pages or runtime signed URLs were inspected.
