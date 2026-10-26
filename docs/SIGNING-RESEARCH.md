# Signing research: macOS precedent and Windows choices

Checked 26 October 2026. This is a research result and proposed direction, not an
implementation or authorization to use credentials, enroll with a provider, sign,
submit, publish, or change another repository. Existing implementation and private
credential values were not changed or read.

## Recommended direction

For macOS, reuse the existing Apple team setup with takedown-specific protected
configuration. Adopt agent-job-queue's separation of build, signing, recoverable
notarization and final acceptance. Keep takedown's simpler app/DMG layout rather
than copying Agent Queue's custom installer stages.

For Windows, prefer **Microsoft Artifact Signing Basic if the publisher is an
eligible Dutch organization**. Its published starting price is **US$9.99/month**
for 5,000 signatures. If publishing as a Dutch individual, this public-trust route
is currently unavailable: Microsoft permits individual applicants only in the US
and Canada. An OSS application to SignPath Foundation or an individual-capable CA
with hardware/cloud key custody becomes the alternative. The publishing identity
has not been assumed. [Eligibility](https://learn.microsoft.com/en-us/azure/artifact-signing/quickstart),
[published USD pricing](https://azure.microsoft.com/en-us/products/artifact-signing),
[Foundation conditions](https://signpath.org/terms.html).

## Verified Apple setup

GitHub's current repository/environment metadata confirms that Hikyo and
agent-job-queue have the same configured Apple team. Both use an `apple-signing`
environment containing these secret names:

- `APPLE_DEVELOPER_ID_P12_BASE64`
- `APPLE_DEVELOPER_ID_P12_PASSWORD`
- `APPLE_NOTARY_KEY_P8`

Their configuration also supplies the signing identity, team ID, notarization key
ID and issuer ID. Metadata confirms configuration exists; it does not reveal or
prove identical certificate/key bytes, current certificate validity, or authorize
copying secrets. Takedown currently uses different input names and a base64 API
key input; an adapter should align the contract deliberately. The existing team
is the precedent, so a new Apple developer enrollment is not the first step.
[Hikyo configuration](https://github.com/Hikyo-Org/Hikyo/blob/4a7beff3272346fb46fec5da058506f1d558282e/docs/release/apple-signing.md),
[Agent Queue configuration](https://github.com/Dunky13/agent-job-queue/blob/b5dac59a30d6365d2edb37e77ba14ff3a2a94416/docs/apple-signing.md).

## How the two macOS pipelines work

The inspected GitHub main commits were Hikyo
`4a7beff3272346fb46fec5da058506f1d558282e` and agent-job-queue
`b5dac59a30d6365d2edb37e77ba14ff3a2a94416`. Exact source snapshots and file hashes
were retained privately. Findings below describe those bytes, not an arbitrary
current branch.

### Hikyo

Builds happen before Apple secrets are loaded. A dedicated signing job imports a
Developer ID Application certificate into a disposable keychain, checks the
identity, signs each architecture's executable with Hardened Runtime and a secure
timestamp, and creates signed DMGs. A separate job receives the notarization key
and verifies the exact signed artifact handoff. It uses `notarytool`, persists
submission receipts and matching Apple logs, requires acceptance, staples and
validates the ticket, verifies the expected Developer ID and runs Gatekeeper
assessment. Checksums/manifests follow the final bytes.
[Signing helper](https://github.com/Hikyo-Org/Hikyo/blob/4a7beff3272346fb46fec5da058506f1d558282e/scripts/release/sign-apple-artifacts.sh),
[release jobs](https://github.com/Hikyo-Org/Hikyo/blob/4a7beff3272346fb46fec5da058506f1d558282e/.github/workflows/release.yml).

Its current `--notarize-only` helper submits again when invoked; preserving
receipts is not a general resume engine. Stable releases require notarization.
Nightlies may use a specifically gated fallback that omits every DMG and declares
unverified notarization. A hosted-runner-only wrapper handles an observed Apple
ticket-network issue; that firewall workaround should not be copied to a local
Mac or enabled without a corresponding need.
[Recovery and fallback contract](https://github.com/Hikyo-Org/Hikyo/blob/4a7beff3272346fb46fec5da058506f1d558282e/docs/release/apple-signing.md).

### Agent Job Queue

Distribution consumes an already-built artifact and verifies its original run,
source identity and digest before credential use. Signing and notarization use
protected main automation and a temporary keychain. The resumable CI path uses a
prebuilt Go client for Apple's official Notary REST API, rather than native
`notarytool` for its direct API-key path. It retains exact signed upload archives,
hashes, stage state and submission IDs in immutable checkpoints. Local one-shot
and keychain-profile paths still use `notarytool`.
[Distribution workflow](https://github.com/Dunky13/agent-job-queue/blob/b5dac59a30d6365d2edb37e77ba14ff3a2a94416/.github/workflows/distribution.yml),
[resume engine](https://github.com/Dunky13/agent-job-queue/blob/b5dac59a30d6365d2edb37e77ba14ff3a2a94416/scripts/macos-resume.sh).

Each create request is attempted once. The client saves the returned submission
ID before uploading to Apple's temporary S3 destination. A retry can poll or
finish an existing stage; it does not blindly resubmit an uncertain request.
Missing evidence, mismatched hashes and unknown outcomes stop recovery. Its
custom installer requires payload, installer and DMG stages. Separate Apple
Silicon and Intel acceptance jobs precede publication and Homebrew verification.
[Submission custody](https://github.com/Dunky13/agent-job-queue/blob/b5dac59a30d6365d2edb37e77ba14ff3a2a94416/internal/notaryapi/submit.go),
[pipeline explanation](https://github.com/Dunky13/agent-job-queue/blob/b5dac59a30d6365d2edb37e77ba14ff3a2a94416/docs/apple-signing.md).

### Hosted evidence checked now

Hikyo run **40171412797**, at source commit
`059902e8ca29ca2d41f7767d0d65a49d9211826f`, reports successful `apple-sign`,
`apple-notarize` and publication jobs. Its source is older than the inspected main
commit, so this is not exact-head qualification of every current main change.
Agent Job Queue distribution **40026971766**, with automation commit
`b5dac59a30d6365d2edb37e77ba14ff3a2a94416`, reports successful
sign/notarize, Apple Silicon/Intel acceptance, publication and Homebrew checks.
These are live workflow results, not a fresh independent download, signature or
Gatekeeper assessment of their published files.
[Hikyo run](https://github.com/Hikyo-Org/Hikyo/actions/runs/40171412797),
[Agent Queue run](https://github.com/Dunky13/agent-job-queue/actions/runs/40026971766).

## What takedown should carry over

1. Build/test the native binary before exposing Apple credentials; consume and
   verify that immutable artifact in the protected signing job. The current
   signed-packaging workflow can build while signing inputs are in its environment.
2. Retain the existing inner-executable then app signing order. Notarize the app
   ZIP, staple/validate the app, recreate the final ZIP, then sign/notarize/staple
   the DMG. A ZIP cannot itself carry a ticket. Takedown has no Agent Queue outer
   custom installer, so it does not need that extra installer stage.
3. Persist artifact hashes, submission IDs and Apple logs outside disposable
   credential scratch. Add explicit resume/finish behavior for the same bytes;
   keep unknown outcomes held. Printing an ID followed by scratch cleanup is not
   equivalent to a durable recovery checkpoint.
4. Keep identity/team/chain verification, keychain restoration, ticket validation
   and a final Gatekeeper assessment. Preserve quarantine during clean-machine
   acceptance. Do not copy another project's nightly fallback as release policy.

These are source-derived recommendations. Apple's workflow supports submission
IDs, log retrieval, a REST API alternative and recreating a ZIP after stapling.
[Apple custom workflow](https://developer.apple.com/documentation/security/customizing-the-notarization-workflow),
[takedown release guide](RELEASING.md).

## Windows choices

Neither inspected project's workflow/scripts contain an Authenticode/SignTool
pipeline to reuse. Hikyo's release signatures/provenance should not be confused
with a Windows publisher signature. Takedown needs its own provider adapter.

### A. Microsoft Artifact Signing, preferred for an eligible organization

Public Trust supports EU organizations, including Dutch organizations. Individual
Public Trust applicants must currently be in the US or Canada. Validation and a
usable Azure subscription/tenant remain prerequisites; no eligibility or business
identity is inferred for the owner.
[Microsoft eligibility](https://learn.microsoft.com/en-us/azure/artifact-signing/quickstart).

Basic lists **US$9.99/month for 5,000 signatures**, then **US$0.005/signature**.
That is US$119.88/year before tax, currency conversion and other CI charges; a
Dutch EUR quote was not established. Managed hardware-backed keys avoid an
exported PFX. GitHub Actions can use OIDC federation and the narrowly scoped
Certificate Profile Signer role instead of a long-lived Azure client secret.
[Pricing](https://azure.microsoft.com/en-us/products/artifact-signing),
[OIDC integration](https://github.com/Azure/artifact-signing-action/blob/main/docs/OIDC.md).

### B. SignPath Foundation, conditional free OSS route

Accepted projects can use free signing with the Foundation's certificate and
publisher identity. Terms require an OSI-approved license, an already released
and maintained project, verifiable reputation, explicit signing roles/approval,
and other published conditions. Takedown currently has no root license declared;
admission is not automatic. Its GitHub connector also needs provenance-preserving
artifact round trips, making the current synchronous NSIS hook a less direct fit.
[Foundation](https://signpath.org/), [terms](https://signpath.org/terms.html),
[GitHub integration](https://docs.signpath.io/trusted-build-systems/github).

### C. Individual-capable CA with token or cloud HSM

A CA offering individual validation is a fallback if Microsoft eligibility or
Foundation admission does not fit. Obtain an applicable Dutch quote and confirm
the exact key-custody/CI integration before purchasing. Sectigo documents
individual applicants, but this research establishes no particular enrollment or
price. A hardware/cloud signer is preferable for unattended CI to a USB token
that requires a physical signing host.
[Individual validation](https://www.sectigo.com/faqs/detail/OV-Code-Signing-Validation-for-Organizations-and-Individuals).

## Necessary Windows adapter changes

The present release adapter expects an exportable PFX. That is not the default
plan for a newly issued public-trust certificate: modern CA key-custody rules
require protected hardware/service storage, and DigiCert explicitly excludes
export of token/HSM/KeyLocker keys. Private/internal trust and legacy exportable
certificates are different cases.
[CA/Browser Forum requirements](https://cabforum.org/working-groups/code-signing/requirements/),
[DigiCert PFX restriction](https://knowledge.digicert.com/general-information/export-a-code-signing-certificate-as-a-pfx-file).

Keep the order **sign application, sign generated uninstaller before embedding,
sign final installer, then hash/archive**. Retain independent Authenticode and
timestamp verification, and add the expected publisher/profile check. A Microsoft
adapter can use SignTool's managed-signing dlib inside NSIS's existing hook.
Its documented timestamp endpoint uses HTTP; the current HTTPS-only timestamp
check needs a deliberate provider-specific policy rather than a blanket relaxation
or an invented HTTPS URL. Foundation requires a reviewed artifact-provenance
staging design instead of treating an arbitrary file upload as trusted output.
[Microsoft integration](https://learn.microsoft.com/en-us/azure/artifact-signing/how-to-signing-integrations),
[NSIS hook](https://nsis.sourceforge.io/Docs/Chapter5.html#uninstfinalize).

Signing does not guarantee zero SmartScreen prompts, and EV is no longer an
automatic bypass. Store/MSIX is a separate distribution route; submitting the
existing EXE/NSIS installer to the Store still requires publisher signing.
[SmartScreen](https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/smartscreen-reputation),
[Store packaging distinction](https://learn.microsoft.com/en-us/windows/apps/package-and-deploy/code-signing-options).

## Decision still needed

Choose the intended Windows publisher identity: eligible organization or
individual. That determines the preferred provider. Licensing, paid enrollment,
provider approval and any production credential provisioning require the owner's
decision. No signing, notarization submission, provider signup, purchase, workflow
dispatch, release publication or changes to Hikyo/Agent Queue occurred here.
