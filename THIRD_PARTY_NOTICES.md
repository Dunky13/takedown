# Third-party notices

This bundle preserves original upstream license, copyright, patent, attribution,
and notice text. Each component retains its own terms; the application's license
does not relicense these components. This inventory is not legal certification.

Scope: the union of the default, cgo-enabled native runtime dependency graphs for
./cmd/takedown on darwin/amd64, darwin/arm64, linux/amd64, and windows/amd64,
using Go 1.26.10. Dependencies reached only through tests or tools are excluded.
Original standalone legal files are included conservatively from each selected
module, so some files may describe unused subcomponents or test assets. Original
source comment notices and embedded font legal name records are also preserved.
Literal native includes shipped within each component are followed recursively.
Vendored OpenGL/GLFW native sources and platform headers are included conservatively;
their inclusion does not assert that each header is linked on every target.
System libraries, SDKs, compiler support libraries, and packaging tools supplied
by build hosts are outside this Go source inventory. Binary distributors must
also account for any separately bundled native libraries or other added payloads.

Regenerate from the repository root after downloading the pinned modules:

    go mod download
    go run -mod=readonly -p 2 scripts/generate-third-party-notices.go

Check without modifying the bundle:

    go run -mod=readonly -p 2 scripts/generate-third-party-notices.go -check

Generation requires Go 1.26.10 and the pinned module cache. The helper verifies
the module cache with go mod verify, prohibits source downloads while generating,
and queries go list -mod=readonly -p 2 -deps -json without -test for all four targets.
Default architecture feature levels, no GOEXPERIMENT, and GOFIPS140=off are used.
No cross compiler or native linking is needed. Targets, components, files, and
deduplicated source notices use stable ordering; output has no timestamp or local
cache paths. Module sums and SHA-256 hashes identify the original source files.
The same pinned inputs produce the same bytes on macOS, Linux, and Windows.

Method references: https://go.dev/ref/mod#go-mod-verify and
https://go.dev/cmd/go/#hdr-List_packages_or_modules.

## Components

- `cloud.google.com/go/auth/oauth2adapt@v0.3.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `cloud.google.com/go/compute/metadata@v0.10.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `fyne.io/fyne/v2@v2.8.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `fyne.io/systray@v1.12.3-0.20260810170012-af4e8e793ec4` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/BurntSushi/toml@v1.6.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/FyshOS/fancyfs@v0.0.1` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/anthonynsimon/bild@v0.14.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/cespare/xxhash/v2@v2.3.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/clipperhouse/uax29/v2@v2.2.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/danieljoos/wincred@v1.2.3` (windows/amd64)

- `github.com/felixge/httpsnoop@v1.1.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/fsnotify/fsnotify@v1.9.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/fyne-io/image@v0.1.1` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/geoffgarside/ber@v1.1.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/go-gl/gl@v0.0.0-20260331235117-4566fea9a276` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/go-logr/logr@v1.4.4` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/go-logr/stdr@v1.2.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/go-text/render@v0.2.1` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/go-text/typesetting@v0.3.5` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/godbus/dbus/v5@v5.2.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/google/s2a-go@v0.1.11` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/google/uuid@v1.6.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/googleapis/enterprise-certificate-proxy@v0.3.23` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/googleapis/gax-go/v2@v2.26.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/hirochachacha/go-smb2@v1.1.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/inconshreveable/mousetrap@v1.1.0` (windows/amd64)

- `github.com/jeandeaual/go-locale@v0.0.0-20250612000132-0ef82f21eade` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/kr/fs@v0.1.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/mattn/go-runewidth@v0.0.24` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/mitchellh/colorstring@v0.0.0-20190213212951-d06e56a500db` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/nicksnyder/go-i18n/v2@v2.5.1` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/pkg/browser@v0.0.0-20240102092130-5ac0b6a4141c` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/pkg/sftp@v1.13.12` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/rivo/uniseg@v0.4.7` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/rymdport/portal@v0.4.2` (linux/amd64)

- `github.com/schollz/progressbar/v3@v3.19.1` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/spf13/cobra@v1.10.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/spf13/pflag@v1.0.9` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/yuin/goldmark@v1.8.2` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `github.com/zalando/go-keyring@v0.2.8` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `go.opentelemetry.io/auto/sdk@v1.2.1` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `go.opentelemetry.io/otel/metric@v1.45.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `go.opentelemetry.io/otel/trace@v1.45.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `go.opentelemetry.io/otel@v1.45.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `go@go1.26.10` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/crypto@v0.58.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/image@v0.47.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/net@v0.61.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/oauth2@v0.37.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/sys@v0.49.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/term@v0.47.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `golang.org/x/text@v0.43.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `google.golang.org/api@v0.308.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `google.golang.org/grpc@v1.84.0` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

- `google.golang.org/protobuf@v1.36.12` (darwin/amd64, darwin/arm64, linux/amd64, windows/amd64)

## cloud.google.com/go/auth/oauth2adapt@v0.3.0

Source: https://pkg.go.dev/cloud.google.com/go/auth/oauth2adapt@v0.3.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:FY8oSZpCYoUNv6QxVODuMjQz4IlSOVeiQtZ08vLPz88=`

### cloud.google.com/go/auth/oauth2adapt@v0.3.0/LICENSE

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- cloud.google.com/go/auth/oauth2adapt@v0.3.0/oauth2adapt.go (source SHA-256: `95e9e512c1f0f8280e8f2f087365d5b296c3c7d20526568d30739d29692b10cc`)

````text
// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a

Source: https://pkg.go.dev/cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:jW1ZLFy9oa5ai8LOI40z6s+DL1cO/L1KSm9hVGnibn8=`

### cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/LICENSE

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/auth.go (source SHA-256: `a68c474a0819ccc52233a8b28d389892fc5bf7c44c3237c50021704725ec4d9e`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/compute.go (source SHA-256: `f486dbed8d44ce39ef6db45592b10d7b5cf734afbf7b36e289f1d9d14b9d0a1b`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/detect.go (source SHA-256: `2621c2091ea562592354cf4f5ce1907abc37c2c9ee7bdde5a158811fe9f3d2c1`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/doc.go (source SHA-256: `4bb020bf380c0f1a44fd22a822be5c15647cc114f7474559deb8c9079c3f0927`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/filetypes.go (source SHA-256: `51f604bfb4acd5185aeb79583d84c0ad16b14c6f1a326476047df0102d0ebf06`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/aws_provider.go (source SHA-256: `ff4a2bfdb63da6b07fe9a15d3d730dfc055ba294ad0a677894d5849ab8f05722`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/executable_provider.go (source SHA-256: `298dff6134ee30b34a8014abb90fc4ed0043e23280530adb3eb83f7f47d0384b`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/externalaccount.go (source SHA-256: `1e00eeb3d41083fd89674ca9eae544707b5dc3e51b0e126cd6ea2e216d1edc1a`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/file_provider.go (source SHA-256: `b0c6dd98de2fbdd61a6dc590b93a2a50b299ff3851b5785ea0424fb6717d1150`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/info.go (source SHA-256: `451136cf77af237bec4cf2402bec4f4c2af2e13f58b0b0b41d80028a92e1a1a6`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/url_provider.go (source SHA-256: `db25b7862c95d21a7e5010ae3e6f74aa80751b4024483f0c076bf31bf5ee59e5`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccountuser/externalaccountuser.go (source SHA-256: `631163a6d137f4557724360f64ec883214c77ccf4521eea468f9fa984a3ed3c7`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/gdch/gdch.go (source SHA-256: `a90b0b200f8210a9eed3b809a0bf06c9edeb223ff903c17529bfd891bf58fcf3`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/impersonate/impersonate.go (source SHA-256: `0cd59ff9dc55d0136392893d682a0bf941d4f02268a143369ad32e4c076e1adf`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/stsexchange/sts_exchange.go (source SHA-256: `1a4d5a8277f8de3c3fb7f86d421e21caf4c8587e0a8acd6f84f484dabfb27962`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/selfsignedjwt.go (source SHA-256: `31c2ca1c40881e9824457a6857e065689862adb5e779634470f8a1930138bdee`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/httptransport/httptransport.go (source SHA-256: `3a6701118b90ea9247f8806d2ceaec6c277355db36816c37dcc3460fd30a0c0d`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/httptransport/transport.go (source SHA-256: `44f38b00c367c1f6319ac1dab1c8433277f6b2d98389a66265e315c8de77cb02`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/credsfile/credsfile.go (source SHA-256: `76d6146298242373b7dc18584f87e14280e43bf543b4c10e6d13cea69219785d`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/credsfile/filetype.go (source SHA-256: `38c392c54b585e0f4533f8ab595e5b7531c63fa47f01242344760cef85324894`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/credsfile/parse.go (source SHA-256: `8c99e346a8b2e1eff6231491cbecf076d760602687e89d4da15e7d046276fa3f`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/internal.go (source SHA-256: `aaa1595df0b4a1c23468fbc8fbb46993f1ddd9c1e9f890e6e8e75809f27ae1bc`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/jwt/jwt.go (source SHA-256: `b05647805bfb39c6426f7061a9e1fab2b777c29a8d7855c8a0587d51241704a4`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/cba.go (source SHA-256: `10f97d025c9592cdc41e01c44d0448e32f11b25e195018e22e5c94d69fc9ef3d`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/cert/default_cert.go (source SHA-256: `ca9bf0d4ffa480ce8ed920f9291a1988eb2e27f79d186536b180b4900703e8ef`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/cert/enterprise_cert.go (source SHA-256: `f2d30ca22c38354777b0f129b1d0e256343e279879caef7f16408ff55480d52e`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/cert/secureconnect_cert.go (source SHA-256: `a5146635ce15302f87dc206b107f1accbc7829e66f35a7d0737e36d4c49b1875`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/s2a.go (source SHA-256: `dbfe9cd182cbf74bce12e06bb0fdb2a5a3df0bf569909bfa5b2f96b2ec627a22`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/transport.go (source SHA-256: `bb13f970c7157cd8c79f465dbb6700c2598a0b274474ab3157402d76a2e2f156`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/threelegged.go (source SHA-256: `89b44c2d264dd59194e94a9e286e2f4bc292e4cbcddb4cd6b691f80ecdcfc7e6`)

````text
// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/programmatic_provider.go (source SHA-256: `7c513b9dd3fb016a416f78f1ceaddd4afaaad4df6499f89cdc384db2d4923e51`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/externalaccount/x509_provider.go (source SHA-256: `1e6b2a7bcb1d48125fcd1c8953a17634fac7e979cee393ae4ec71e0a037d0b62`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/cert/workload_cert.go (source SHA-256: `9abfb3b4a9a3168c38374f2efd6bad6dc6f6001d1892162f0d01c90e9a64a9f2`)

````text
// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/credentials/internal/impersonate/idtoken.go (source SHA-256: `9b61b386c9907f8f070b279fc64d9b633209eb1f086ca0222fe24f8581850377`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/regionalaccessboundary/external_accounts_config_providers.go (source SHA-256: `5a8b2f2dc3433aa4765884e1fa203cc125e36553999f84aedee0646de933d604`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/retry/retry.go (source SHA-256: `ca098e8222c08f3bec247d62ee9339f383f7e27cefce3257eafc9b3d9398c48b`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/transport/headers/headers.go (source SHA-256: `c01953e9be8bad7a72a84aef2ec762d3d37d3dd0eee538709bc9ecbcca4e94ee`)

````text
// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/regionalaccessboundary/regional_access_boundary.go (source SHA-256: `462db94f1526399194f111095de238985d3b68ebd423da751596f8fe1c82156c`)
- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/retry/retry_linux.go (source SHA-256: `86db342d02267a61acb1f63b85bfc5a151b21d9e1460fc09c4ae28fc7b38100f`)

````text
// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- cloud.google.com/go/auth@v0.24.1-0.20261001053825-dbc26066f70a/internal/version.go (source SHA-256: `24598ae84019253d78729668298c64af37a1e451c8b74bd779d10c6a05216f04`)

````text
// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## cloud.google.com/go/compute/metadata@v0.10.0

Source: https://pkg.go.dev/cloud.google.com/go/compute/metadata@v0.10.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:pyKMUQSwchgkIBBJGdILqQbs/BNJXqwSA7Ej6LAvvtY=`

### cloud.google.com/go/compute/metadata@v0.10.0/LICENSE

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- cloud.google.com/go/compute/metadata@v0.10.0/metadata.go (source SHA-256: `cc964271dba9aa8928c6223216d6ea9d188b9ecdee963e6d10b54a189d4cbc2d`)

````text
// Copyright 2014 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- cloud.google.com/go/compute/metadata@v0.10.0/retry.go (source SHA-256: `4d87e7da851bfa31bb83361f6dc47e4e8c39d2c93a5b85ced68acbfcbb1f73e8`)
- cloud.google.com/go/compute/metadata@v0.10.0/retry_linux.go (source SHA-256: `965e03771ec5b835462260af5c7f3571e41b0843cf24d01b093bab3204853880`)

````text
// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- cloud.google.com/go/compute/metadata@v0.10.0/log.go (source SHA-256: `ce359ee6898550a082cbec40acbaa9b12a60d1f463731f019c0ce67293bb5b0e`)
- cloud.google.com/go/compute/metadata@v0.10.0/syscheck.go (source SHA-256: `38242101ed8e8c5dbf4a976681dcecf8115b481cb7332e760c12b51415c9239f`)
- cloud.google.com/go/compute/metadata@v0.10.0/syscheck_linux.go (source SHA-256: `d3b02acb20715e6ef6dff2870c6400811f00352c63dec1da4412f901535194f0`)
- cloud.google.com/go/compute/metadata@v0.10.0/syscheck_windows.go (source SHA-256: `535ee98e1c8443502eeccbb62e69e8c972ae0fa74bd4dbe72600500445790431`)

````text
// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## fyne.io/fyne/v2@v2.8.2

Source: https://pkg.go.dev/fyne.io/fyne/v2@v2.8.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:gYNJj+ZEPd2wF3M8SlF5iL6febpZkZ1csheRFkuEfHU=`

### fyne.io/fyne/v2@v2.8.2/AUTHORS

SHA-256: `477326325c45f37a21117797cd40b1157c6090f9e4a2b5fb839d6bbb9b595453`

````text
Andy Williams <andy@andy.xyz>
Steve OConnor <steveoc64@gmail.com>
Luca Corbo <lu.corbo@gmail.com>
Paul Hovey <paul@paulhovey.org>
Charles Corbett <nafredy@gmail.com>
Tilo Prütz <tilo@pruetz.net>
Stephen Houston <smhouston88@gmail.com>
Storm Hess <stormhess@gloryskulls.com>
Stuart Scott <stuart.murray.scott@gmail.com>
Jacob Alzén <jacalz@tutanota.com>
Charles A. Daniels <charles@cdaniels.net>
Pablo Fuentes <f.pablo1@hotmail.com>
Changkun Ou <hi@changkun.de>
Cedric Bail
Drew Weymouth
Simon Dassow
Max Schillinger
````

### fyne.io/fyne/v2@v2.8.2/LICENSE

SHA-256: `68a1eafb90581708e62e1d269e7f6e92e6b0b5cdae93f45ede832d42912c7014`

````text
BSD 3-Clause License

Copyright (C) 2018 Fyne.io developers (see AUTHORS)
All rights reserved.


Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of Fyne.io nor the names of its contributors may be
      used to endorse or promote products derived from this software without
      specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

````

### fyne.io/fyne/v2@v2.8.2/theme/font/LICENSE.txt

SHA-256: `c4910f5c1bcf2f7326f9fc3d33b7fc4eed7bf294374e3a17aaee64ceff0fb85b`

````text
—————————————————————————————-
SIL OPEN FONT LICENSE Version 1.1 - 26 February 2007
—————————————————————————————-

PREAMBLE
The goals of the Open Font License (OFL) are to stimulate worldwide development of collaborative font projects, to support the font creation efforts of academic and linguistic communities, and to provide a free and open framework in which fonts may be shared and improved in partnership with others.

The OFL allows the licensed fonts to be used, studied, modified and redistributed freely as long as they are not sold by themselves. The fonts, including any derivative works, can be bundled, embedded, redistributed and/or sold with any software provided that any reserved names are not used by derivative works. The fonts and derivatives, however, cannot be released under any other type of license. The requirement for fonts to remain under this license does not apply to any document created using the fonts or their derivatives.

DEFINITIONS
“Font Software” refers to the set of files released by the Copyright Holder(s) under this license and clearly marked as such. This may include source files, build scripts and documentation.

“Reserved Font Name” refers to any names specified as such after the copyright statement(s).

“Original Version” refers to the collection of Font Software components as distributed by the Copyright Holder(s).

“Modified Version” refers to any derivative made by adding to, deleting, or substituting—in part or in whole—any of the components of the Original Version, by changing formats or by porting the Font Software to a new environment.

“Author” refers to any designer, engineer, programmer, technical writer or other person who contributed to the Font Software.

PERMISSION & CONDITIONS
Permission is hereby granted, free of charge, to any person obtaining a copy of the Font Software, to use, study, copy, merge, embed, modify, redistribute, and sell modified and unmodified copies of the Font Software, subject to the following conditions:

1) Neither the Font Software nor any of its individual components, in Original or Modified Versions, may be sold by itself.

2) Original or Modified Versions of the Font Software may be bundled, redistributed and/or sold with any software, provided that each copy contains the above copyright notice and this license. These can be included either as stand-alone text files, human-readable headers or in the appropriate machine-readable metadata fields within text or binary files as long as those fields can be easily viewed by the user.

3) No Modified Version of the Font Software may use the Reserved Font Name(s) unless explicit written permission is granted by the corresponding Copyright Holder. This restriction only applies to the primary font name as presented to the users.

4) The name(s) of the Copyright Holder(s) or the Author(s) of the Font Software shall not be used to promote, endorse or advertise any Modified Version, except to acknowledge the contribution(s) of the Copyright Holder(s) and the Author(s) or with their explicit written permission.

5) The Font Software, modified or unmodified, in part or in whole, must be distributed entirely under this license, and must not be distributed under any other license. The requirement for fonts to remain under this license does not apply to any document created using the Font Software.

TERMINATION
This license becomes null and void if any of the above conditions are not met.

DISCLAIMER
THE FONT SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM OTHER DEALINGS IN THE FONT SOFTWARE.
````

### fyne.io/fyne/v2@v2.8.2/theme/font/LICENSE_DejaVu-Powerline.txt

SHA-256: `a30f2466a95cd1eaeff6041692f43469edcca052d002e877015a819c6e31c194`

````text
Fonts are (c) Bitstream (see below). DejaVu changes are in public domain.
Glyphs imported from Arev fonts are (c) Tavmjong Bah (see below)

Bitstream Vera Fonts Copyright
------------------------------

Copyright (c) 2003 by Bitstream, Inc. All Rights Reserved. Bitstream Vera is
a trademark of Bitstream, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of the fonts accompanying this license ("Fonts") and associated
documentation files (the "Font Software"), to reproduce and distribute the
Font Software, including without limitation the rights to use, copy, merge,
publish, distribute, and/or sell copies of the Font Software, and to permit
persons to whom the Font Software is furnished to do so, subject to the
following conditions:

The above copyright and trademark notices and this permission notice shall
be included in all copies of one or more of the Font Software typefaces.

The Font Software may be modified, altered, or added to, and in particular
the designs of glyphs or characters in the Fonts may be modified and
additional glyphs or characters may be added to the Fonts, only if the fonts
are renamed to names not containing either the words "Bitstream" or the word
"Vera".

This License becomes null and void to the extent applicable to Fonts or Font
Software that has been modified and is distributed under the "Bitstream
Vera" names.

The Font Software may be sold as part of a larger software package but no
copy of one or more of the Font Software typefaces may be sold by itself.

THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS
OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT OF COPYRIGHT, PATENT,
TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL BITSTREAM OR THE GNOME
FOUNDATION BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, INCLUDING
ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL DAMAGES,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF
THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM OTHER DEALINGS IN THE
FONT SOFTWARE.

Except as contained in this notice, the names of Gnome, the Gnome
Foundation, and Bitstream Inc., shall not be used in advertising or
otherwise to promote the sale, use or other dealings in this Font Software
without prior written authorization from the Gnome Foundation or Bitstream
Inc., respectively. For further information, contact: fonts at gnome dot
org. 

Arev Fonts Copyright
------------------------------

Copyright (c) 2006 by Tavmjong Bah. All Rights Reserved.

Permission is hereby granted, free of charge, to any person obtaining
a copy of the fonts accompanying this license ("Fonts") and
associated documentation files (the "Font Software"), to reproduce
and distribute the modifications to the Bitstream Vera Font Software,
including without limitation the rights to use, copy, merge, publish,
distribute, and/or sell copies of the Font Software, and to permit
persons to whom the Font Software is furnished to do so, subject to
the following conditions:

The above copyright and trademark notices and this permission notice
shall be included in all copies of one or more of the Font Software
typefaces.

The Font Software may be modified, altered, or added to, and in
particular the designs of glyphs or characters in the Fonts may be
modified and additional glyphs or characters may be added to the
Fonts, only if the fonts are renamed to names not containing either
the words "Tavmjong Bah" or the word "Arev".

This License becomes null and void to the extent applicable to Fonts
or Font Software that has been modified and is distributed under the 
"Tavmjong Bah Arev" names.

The Font Software may be sold as part of a larger software package but
no copy of one or more of the Font Software typefaces may be sold by
itself.

THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT
OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL
TAVMJONG BAH BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL
DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM
OTHER DEALINGS IN THE FONT SOFTWARE.

Except as contained in this notice, the name of Tavmjong Bah shall not
be used in advertising or otherwise to promote the sale, use or other
dealings in this Font Software without prior written authorization
from Tavmjong Bah. For further information, contact: tavmjong @ free
. fr.
````

### fyne.io/fyne/v2@v2.8.2/theme/font/LICENSE_EmojiOneColor.txt

SHA-256: `a79e11d824be724d362c8c143b4130889c66cf2ef47f16e5ad83bd796dd18183`

````text
MIT License

Copyright (c) 2016 Adobe Systems Incorporated

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````

### fyne.io/fyne/v2@v2.8.2/theme/font/LICENSE_Inter.txt

SHA-256: `4d7d9c95e7d7f2f0ebf76d5e0b344826b74e903a34028ee18ab54bb639e45906`

````text
Copyright (c) 2016-2020 The Inter Project Authors.
"Inter" is trademark of Rasmus Andersson.
https://github.com/rsms/inter

This Font Software is licensed under the SIL Open Font License, Version 1.1.
This license is copied below, and is also available with a FAQ at:
http://scripts.sil.org/OFL

-----------------------------------------------------------
SIL OPEN FONT LICENSE Version 1.1 - 26 February 2007
-----------------------------------------------------------

PREAMBLE
The goals of the Open Font License (OFL) are to stimulate worldwide
development of collaborative font projects, to support the font creation
efforts of academic and linguistic communities, and to provide a free and
open framework in which fonts may be shared and improved in partnership
with others.

The OFL allows the licensed fonts to be used, studied, modified and
redistributed freely as long as they are not sold by themselves. The
fonts, including any derivative works, can be bundled, embedded,
redistributed and/or sold with any software provided that any reserved
names are not used by derivative works. The fonts and derivatives,
however, cannot be released under any other type of license. The
requirement for fonts to remain under this license does not apply
to any document created using the fonts or their derivatives.

DEFINITIONS
"Font Software" refers to the set of files released by the Copyright
Holder(s) under this license and clearly marked as such. This may
include source files, build scripts and documentation.

"Reserved Font Name" refers to any names specified as such after the
copyright statement(s).

"Original Version" refers to the collection of Font Software components as
distributed by the Copyright Holder(s).

"Modified Version" refers to any derivative made by adding to, deleting,
or substituting -- in part or in whole -- any of the components of the
Original Version, by changing formats or by porting the Font Software to a
new environment.

"Author" refers to any designer, engineer, programmer, technical
writer or other person who contributed to the Font Software.

PERMISSION AND CONDITIONS
Permission is hereby granted, free of charge, to any person obtaining
a copy of the Font Software, to use, study, copy, merge, embed, modify,
redistribute, and sell modified and unmodified copies of the Font
Software, subject to the following conditions:

1) Neither the Font Software nor any of its individual components,
in Original or Modified Versions, may be sold by itself.

2) Original or Modified Versions of the Font Software may be bundled,
redistributed and/or sold with any software, provided that each copy
contains the above copyright notice and this license. These can be
included either as stand-alone text files, human-readable headers or
in the appropriate machine-readable metadata fields within text or
binary files as long as those fields can be easily viewed by the user.

3) No Modified Version of the Font Software may use the Reserved Font
Name(s) unless explicit written permission is granted by the corresponding
Copyright Holder. This restriction only applies to the primary font name as
presented to the users.

4) The name(s) of the Copyright Holder(s) or the Author(s) of the Font
Software shall not be used to promote, endorse or advertise any
Modified Version, except to acknowledge the contribution(s) of the
Copyright Holder(s) and the Author(s) or with their explicit written
permission.

5) The Font Software, modified or unmodified, in part or in whole,
must be distributed entirely under this license, and must not be
distributed under any other license. The requirement for fonts to
remain under this license does not apply to any document created
using the Font Software.

TERMINATION
This license becomes null and void if any of the above conditions are
not met.

DISCLAIMER
THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT
OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL THE
COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL
DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM
OTHER DEALINGS IN THE FONT SOFTWARE.
````

### fyne.io/fyne/v2@v2.8.2/theme/font/DejaVuSansMono-Powerline.ttf (font legal name records)

Source SHA-256: `af1868d514095dc61f090a876e4413a4ebadf4c81e8748f9fcbb7c69eac95853`

````text
Copyright:
Copyright (c) 2003 by Bitstream, Inc. All Rights Reserved.
DejaVu changes are in public domain


License:
Fonts are (c) Bitstream (see below). DejaVu changes are in public domain.

Bitstream Vera Fonts Copyright
------------------------------

Copyright (c) 2003 by Bitstream, Inc. All Rights Reserved. Bitstream Vera is a trademark of Bitstream, Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy of the fonts accompanying this license ("Fonts") and associated documentation files (the "Font Software"), to reproduce and distribute the Font Software, including without limitation the rights to use, copy, merge, publish, distribute, and/or sell copies of the Font Software, and to permit persons to whom the Font Software is furnished to do so, subject to the following conditions:

The above copyright and trademark notices and this permission notice shall be included in all copies of one or more of the Font Software typefaces.

The Font Software may be modified, altered, or added to, and in particular the designs of glyphs or characters in the Fonts may be modified and additional glyphs or  or characters may be added to the Fonts, only if the fonts are renamed to names not containing either the words "Bitstream" or the word "Vera".

This License becomes null and void to the extent applicable to Fonts or Font Software that has been modified and is distributed under the "Bitstream Vera" names.

The Font Software may be sold as part of a larger software package but no copy of one or more of the Font Software typefaces may be sold by itself.

THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL BITSTREAM OR THE GNOME FOUNDATION BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM OTHER DEALINGS IN THE FONT SOFTWARE.

Except as contained in this notice, the names of Gnome, the Gnome Foundation, and Bitstream Inc., shall not be used in advertising or otherwise to promote the sale, use or other dealings in this Font Software without prior written authorization from the Gnome Foundation or Bitstream Inc., respectively. For further information, contact: fonts at gnome dot org. 


License URL:
http://dejavu.sourceforge.net/wiki/index.php/License

````

### fyne.io/fyne/v2@v2.8.2/theme/font/EmojiOneColor.otf (font legal name records)

Source SHA-256: `e9ec7cee76c09ac8d70b1e54ad3f44ae2d472a3441cec3ca4f34ef00b2a0b377`

````text
Copyright:
Copyright 2016 Adobe Systems Incorporated

License:
EmojiOne's graphics are free to use for any project, commercial or personal, under a free culture Creative Commons License (CC-BY 4.0). Proper attribution (link back) is required for the rights to use the emoji in commercial projects.

License URL:
http://emojione.com/licensing

````

### fyne.io/fyne/v2@v2.8.2/theme/font/InterSymbols-Regular.ttf (font legal name records)

Source SHA-256: `f7f79743e9ef5a31db29508ce5a29fa69fc8e50ea837143ebae0fe002954a77a`

````text
Copyright:
Copyright © 2020 The Inter Project Authors

````

### fyne.io/fyne/v2@v2.8.2/theme/font/NotoSans-Bold.ttf (font legal name records)

Source SHA-256: `7c15ac396d2ce6bc33a3b4efacdbd322c9e46376599a725c6f790d8036052cab`

````text
Copyright:
Copyright 2015 Google Inc. All Rights Reserved.

Trademark:
Noto is a trademark of Google Inc.

License:
This Font Software is licensed under the SIL Open Font License, Version 1.1. This Font Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the SIL Open Font License for the specific language, permissions and limitations governing your use of this Font Software.

License URL:
http://scripts.sil.org/OFL

````

### fyne.io/fyne/v2@v2.8.2/theme/font/NotoSans-BoldItalic.ttf (font legal name records)

Source SHA-256: `871c3129a8ba4dcd6b37949c9931f1087ff0434731e6db15798f00fc98808480`

````text
Copyright:
Copyright 2015 Google Inc. All Rights Reserved.

Trademark:
Noto is a trademark of Google Inc.

License:
This Font Software is licensed under the SIL Open Font License, Version 1.1. This Font Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the SIL Open Font License for the specific language, permissions and limitations governing your use of this Font Software.

License URL:
http://scripts.sil.org/OFL

````

### fyne.io/fyne/v2@v2.8.2/theme/font/NotoSans-Italic.ttf (font legal name records)

Source SHA-256: `35600501daadd3e9030ee6ed8098fd3c7f31bfd898ca58da8edf793a3231c48b`

````text
Copyright:
Copyright 2015 Google Inc. All Rights Reserved.

Trademark:
Noto is a trademark of Google Inc.

License:
This Font Software is licensed under the SIL Open Font License, Version 1.1. This Font Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the SIL Open Font License for the specific language, permissions and limitations governing your use of this Font Software.

License URL:
http://scripts.sil.org/OFL

````

### fyne.io/fyne/v2@v2.8.2/theme/font/NotoSans-Regular.ttf (font legal name records)

Source SHA-256: `6b04c8dd65af6b73eb4279472ed1580b29102d6496a377340e80a40cdb3b22c9`

````text
Copyright:
Copyright 2015 Google Inc. All Rights Reserved.

Trademark:
Noto is a trademark of Google Inc.

License:
This Font Software is licensed under the SIL Open Font License, Version 1.1. This Font Software is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the SIL Open Font License for the specific language, permissions and limitations governing your use of this Font Software.

License URL:
http://scripts.sil.org/OFL

````


## fyne.io/systray@v1.12.3-0.20260810170012-af4e8e793ec4

Source: https://pkg.go.dev/fyne.io/systray@v1.12.3-0.20260810170012-af4e8e793ec4

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:149/+Wa5EsLLXfyj2pdTmvnQf2VIlgCIwSjcCTHYhIo=`

### fyne.io/systray@v1.12.3-0.20260810170012-af4e8e793ec4/LICENSE

SHA-256: `a106bc9f4d5d7d91d5173e86d0864fcf061c33a30cc237710b8f4a05cf127021`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright 2014 Brave New Software Project, Inc.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````


## github.com/BurntSushi/toml@v1.6.0

Source: https://pkg.go.dev/github.com/BurntSushi/toml@v1.6.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:dRaEfpa2VI55EwlIW72hMRHdWouJeRF7TPYhI+AUQjk=`

### github.com/BurntSushi/toml@v1.6.0/COPYING

SHA-256: `d21cb1c60785d6d3a84a7059323ccafc45c645b1bbda281c76a62d66ad2d7dc3`

````text
The MIT License (MIT)

Copyright (c) 2013 TOML authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

### github.com/BurntSushi/toml@v1.6.0/cmd/toml-test-decoder/COPYING

SHA-256: `d21cb1c60785d6d3a84a7059323ccafc45c645b1bbda281c76a62d66ad2d7dc3`

````text
The MIT License (MIT)

Copyright (c) 2013 TOML authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

### github.com/BurntSushi/toml@v1.6.0/cmd/toml-test-encoder/COPYING

SHA-256: `d21cb1c60785d6d3a84a7059323ccafc45c645b1bbda281c76a62d66ad2d7dc3`

````text
The MIT License (MIT)

Copyright (c) 2013 TOML authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

### github.com/BurntSushi/toml@v1.6.0/cmd/tomlv/COPYING

SHA-256: `d21cb1c60785d6d3a84a7059323ccafc45c645b1bbda281c76a62d66ad2d7dc3`

````text
The MIT License (MIT)

Copyright (c) 2013 TOML authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

### Original source notice

- github.com/BurntSushi/toml@v1.6.0/type_fields.go (source SHA-256: `b557f6362f1ffee1f5e07dda1b7be72c7579be465d0916e70e7a327129c86000`)

````text
// Struct field handling is adapted from code in encoding/json:
//
// Copyright 2010 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the Go distribution.
````


## github.com/FyshOS/fancyfs@v0.0.1

Source: https://pkg.go.dev/github.com/FyshOS/fancyfs@v0.0.1

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:kgvm7VvwOMLkYTqSflplp62SlMVWQ2uAoHw9CXwXHYg=`

### github.com/FyshOS/fancyfs@v0.0.1/AUTHORS

SHA-256: `efacd7afd0856fee05c335bd2d1b183c289aa013be54295b64ded0d56071461e`

````text
Andy Williams <andy@andy.xyz>
````

### github.com/FyshOS/fancyfs@v0.0.1/LICENSE

SHA-256: `6df86a72d288d232476070e4534f35fcec1dd830ad303d76c7f577e271a369d2`

````text
Copyright (C) 2025 FyshOS developers (see AUTHORS)
All rights reserved.


Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name of Fyne.io nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

````


## github.com/anthonynsimon/bild@v0.14.0

Source: https://pkg.go.dev/github.com/anthonynsimon/bild@v0.14.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:IFRkmKdNdqmexXHfEU7rPlAmdUZ8BDZEGtGHDnGWync=`

### github.com/anthonynsimon/bild@v0.14.0/LICENSE

SHA-256: `637516cf71f2b3eda480c6cb39ee62297d8748acc95d8b4c63200d04e18793ea`

````text
MIT License

Copyright (c) 2021 Anthony Najjar Simon

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/cespare/xxhash/v2@v2.3.0

Source: https://pkg.go.dev/github.com/cespare/xxhash/v2@v2.3.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=`

### github.com/cespare/xxhash/v2@v2.3.0/LICENSE.txt

SHA-256: `f566a9f97bacdaf00d9f21dd991e81dc11201c4e016c86b470799429a1c9a79c`

````text
Copyright (c) 2016 Caleb Spare

MIT License

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
````


## github.com/clipperhouse/uax29/v2@v2.2.0

Source: https://pkg.go.dev/github.com/clipperhouse/uax29/v2@v2.2.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:ChwIKnQN3kcZteTXMgb1wztSgaU+ZemkgWdohwgs8tY=`

### github.com/clipperhouse/uax29/v2@v2.2.0/LICENSE

SHA-256: `07e9e3daaccc707c9b7d535ed1405a65d37919fdd42197fe55e68820e974688c`

````text
MIT License

Copyright (c) 2020 Matt Sherman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/danieljoos/wincred@v1.2.3

Source: https://pkg.go.dev/github.com/danieljoos/wincred@v1.2.3

Targets: windows/amd64

Module checksum: `h1:v7dZC2x32Ut3nEfRH+vhoZGvN72+dQ/snVXo/vMFLdQ=`

### github.com/danieljoos/wincred@v1.2.3/LICENSE

SHA-256: `1702a813fa31858c0d6c7788216fa3dd202472265716998024654c252fc7bb41`

````text
The MIT License (MIT)

Copyright (c) 2014 Daniel Joos

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/felixge/httpsnoop@v1.1.0

Source: https://pkg.go.dev/github.com/felixge/httpsnoop@v1.1.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:3YtUj32ZZkqZtt3sZZsClsymw/QDuVfpNhoA31zeORc=`

### github.com/felixge/httpsnoop@v1.1.0/LICENSE.txt

SHA-256: `8108e14b5cb2b6bcd51583d3849ef0c2b5c4de61758256c20308bc70d1775dbb`

````text
Copyright (c) 2016 Felix Geisendörfer (felix@debuggable.com)

 Permission is hereby granted, free of charge, to any person obtaining a copy
 of this software and associated documentation files (the "Software"), to deal
 in the Software without restriction, including without limitation the rights
 to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 copies of the Software, and to permit persons to whom the Software is
 furnished to do so, subject to the following conditions:

 The above copyright notice and this permission notice shall be included in
 all copies or substantial portions of the Software.

 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 THE SOFTWARE.
````


## github.com/fsnotify/fsnotify@v1.9.0

Source: https://pkg.go.dev/github.com/fsnotify/fsnotify@v1.9.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:2Ml+OJNzbYCTzsxtv8vKSFD9PbJjmhYF14k/jKC7S9k=`

### github.com/fsnotify/fsnotify@v1.9.0/LICENSE

SHA-256: `2dfbe6d2b5eb18b55148f3805d88648a09e0d17e38428123259afb08c8996e37`

````text
Copyright © 2012 The Go Authors. All rights reserved.
Copyright © fsnotify Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without modification,
are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above copyright notice, this
  list of conditions and the following disclaimer in the documentation and/or
  other materials provided with the distribution.
* Neither the name of Google Inc. nor the names of its contributors may be used
  to endorse or promote products derived from this software without specific
  prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/fyne-io/image@v0.1.1

Source: https://pkg.go.dev/github.com/fyne-io/image@v0.1.1

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:WH0z4H7qfvNUw5l4p3bC1q70sa5+YWVt6HCj7y4VNyA=`

### github.com/fyne-io/image@v0.1.1/LICENSE

SHA-256: `b9c31d8f1a43d7260789099ef59cc17a8b72aa3fccd1479bd516ec1b924637c3`

````text
BSD 3-Clause License

Copyright (c) 2022, Fyne.io
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### github.com/fyne-io/image@v0.1.1/ico/LICENSE-reader

SHA-256: `b40930bbcf80744c86c46a12bc9da056641d722716c378f5659b9e555ef833e1`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "{}"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright {yyyy} {name of copyright owner}

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### github.com/fyne-io/image@v0.1.1/ico/LICENSE-writer

SHA-256: `dd26a7abddd02e2d0aba97805b31f248ef7835d9e10da289b22e3b8ab78b324d`

````text
Copyright (c) 2012 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9

Source: https://pkg.go.dev/github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:owsdpq6a2KhaAA/hoHGQmFM7kOgKo3srXZVGnbEUzr0=`

### github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/LICENSE

SHA-256: `cda28abc46a25d34ad876ec119ad3c29edbb8634d5fef46ccf2e1d63acd9ca34`

````text
BSD 3-Clause License

Copyright (c) 2018, Steven R Wiley
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/testdata/LICENSE

SHA-256: `c654732ce4d4761102c29a97006abf63f2c2294c1558dd73d6acca5a327840f8`

````text
Test icons in the testdata folder are from either www.freepik.com or www.flaticon.com. They request any wedsite using these icons place the following in thei web pages:

<div>Icons made by <a href="http://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>

<div>Icons made by <a href="https://www.flaticon.com/authors/flat-icons" title="Flat Icons">Flat Icons</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
````

### Original source notice

- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/definitions.go (source SHA-256: `7df00b0b4bb1151a2afb1ca2c69b5fdc4d3a1db1dcb9f8ca37b284ee820b74c8`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/draw.go (source SHA-256: `52499789e49bd75bc6c940ad5b51a8640534cdc0ee2db35cc92e2a90a2c715eb`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/icon_cursor.go (source SHA-256: `748eb1ad3964696d2ad46198af0890cc9c3874314bfeeaa962ef2e74e4982424`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/path_cursor.go (source SHA-256: `490a50ba5171827068e158491258188ed906dd7975794a279f40eae45d22e940`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/path_style.go (source SHA-256: `551259dc15aed7c44e3bea6abe10da7317b353c7b03b0b2e8ac3db0f1b220219`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/public.go (source SHA-256: `a2face8c9be4670a49542a3a6d412aa343f99ad87e5c92de266ed3f953675100`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/svg_icon.go (source SHA-256: `087831bf805003e9793c101680988cc6447303423b957880f0d36053a8f207e5`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/svg_path.go (source SHA-256: `934ec3bfb250e4377eacba11e8325dc83e7396fb7f0992a8a1a9e2a2b75de856`)
- github.com/fyne-io/oksvg@v0.2.1-0.20260918172519-a9af55fa95d9/utils.go (source SHA-256: `c8f617ac250c5aa491b53493cb2f83e888d532c910658634149c2fdc3318986a`)

````text
// Copyright 2017 The oksvg Authors. All rights reserved.
// created: 2/12/2017 by S.R.Wiley
//
// utils.go implements translation of an SVG2.0 path into a rasterx Path.
````


## github.com/geoffgarside/ber@v1.1.0

Source: https://pkg.go.dev/github.com/geoffgarside/ber@v1.1.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:qTmFG4jJbwiSzSXoNJeHcOprVzZ8Ulde2Rrrifu5U9w=`

### github.com/geoffgarside/ber@v1.1.0/LICENSE

SHA-256: `2d36597f7117c38b006835ae7f537487207d8ec407aa9d9980794b2030cbc067`

````text
Copyright (c) 2009 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/geoffgarside/ber@v1.1.0/ber.go (source SHA-256: `cb8884b8ab46106f86779b2a74848138709e79b715dfc6398e3afe0b5a11a625`)
- github.com/geoffgarside/ber@v1.1.0/common.go (source SHA-256: `def6a3d49e63e2766a40a8505581222b962f49d72150b20492ae11c372facc5e`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## github.com/go-gl/gl@v0.0.0-20260331235117-4566fea9a276

Source: https://pkg.go.dev/github.com/go-gl/gl@v0.0.0-20260331235117-4566fea9a276

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:IO5P06Pcj9K04d+l4nrf3c2U56+dAotIFG6u4P1wAHI=`

### github.com/go-gl/gl@v0.0.0-20260331235117-4566fea9a276/LICENSE

SHA-256: `035b8c4d5796e3465c668aa0d6df233f644a7aecdeeb9016dcec1b2875f70c6b`

````text
The MIT License (MIT)

Copyright (c) 2014 Eric Woroshow

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### Original source notice

- github.com/go-gl/gl@v0.0.0-20260331235117-4566fea9a276/v2.1/gl/package.go (source SHA-256: `7c37e8baa148dfe0b64dc0987149fd65dd9041fc5c7437fc36e2cf84bf0d5642`)

````text
// Copyright (c) 2010 Khronos Group.
// This material may be distributed subject to the terms and conditions
// set forth in the Open Publication License, v 1.0, 8 June 1999.
// http://opencontent.org/openpub/.
//
// Copyright (c) 1991-2006 Silicon Graphics, Inc.
// This document is licensed under the SGI Free Software B License.
// For details, see http://oss.sgi.com/projects/FreeB.
````


## github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2

Source: https://pkg.go.dev/github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:ez76TQInhkDIMobWbOFcPJw6FeB7coXRAVZpKwF8NBg=`

### github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/LICENSE

SHA-256: `72f474049924563f43e5a5cc8e7e228a4fc8eb58e3dd7383bd3552c0725f566c`

````text
Copyright (c) 2012 The glfw3-go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/LICENSE.md

SHA-256: `149704059b5d0bf551637e50042dd4de9c2cae921021f6636298911e3a5f9462`

````text
Copyright (c) 2002-2006 Marcus Geelnard

Copyright (c) 2006-2019 Camilla Löwy

This software is provided 'as-is', without any express or implied
warranty. In no event will the authors be held liable for any damages
arising from the use of this software.

Permission is granted to anyone to use this software for any purpose,
including commercial applications, and to alter it and redistribute it
freely, subject to the following restrictions:

1. The origin of this software must not be misrepresented; you must not
   claim that you wrote the original software. If you use this software
   in a product, an acknowledgment in the product documentation would
   be appreciated but is not required.

2. Altered source versions must be plainly marked as such, and must not
   be misrepresented as being the original software.

3. This notice may not be removed or altered from any source
   distribution.

````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/deps/mingw/dinput.h (source SHA-256: `72316ef3309b6a9b60d16f064a84fc5af6efd84546d69ac557b117bc398e3aab`)

````text
/*
 * Copyright (C) the Wine project
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301, USA
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/wayland-client-protocol-code.h (source SHA-256: `fcba248192c17bd5334129c92c49c88554b060c45196beaf5de9f8c5a552febc`)

````text
/*
 * Copyright © 2008-2011 Kristian Høgsberg
 * Copyright © 2010-2011 Intel Corporation
 * Copyright © 2012-2013 Collabora, Ltd.
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the
 * next paragraph) shall be included in all copies or substantial
 * portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT.  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/xdg-shell-client-protocol-code.h (source SHA-256: `e3507c8869d70682af7580ff574894c96ba4eaa5804a61c3335a5bb38564ff8c`)

````text
/*
 * Copyright © 2008-2013 Kristian Høgsberg
 * Copyright © 2013      Rafael Antognolli
 * Copyright © 2013      Jasper St. Pierre
 * Copyright © 2010-2013 Intel Corporation
 * Copyright © 2015-2017 Samsung Electronics Co., Ltd
 * Copyright © 2015-2017 Red Hat Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/viewporter-client-protocol-code.h (source SHA-256: `5164cd057303e511d2a202eeaec9074997e47ddfe6500a5e26354f85e86417cf`)

````text
/*
 * Copyright © 2013-2016 Collabora, Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/pointer-constraints-unstable-v1-client-protocol-code.h (source SHA-256: `376b72b7ac9d7c3f48279c5f9e1fd778997007d6e37e014e6febf045b853e368`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/relative-pointer-unstable-v1-client-protocol-code.h (source SHA-256: `e91835b11afe3e1759e7c946f1dd9dc7e6fc3f1de89cd3455e43e7df66531429`)

````text
/*
 * Copyright © 2014      Jonas Ådahl
 * Copyright © 2015      Red Hat Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/idle-inhibit-unstable-v1-client-protocol-code.h (source SHA-256: `3166872808c85879adaafd6dcdcce0113bbe491eca572aee35b011ea946add75`)

````text
/*
 * Copyright © 2015 Samsung Electronics Co., Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/xdg-decoration-unstable-v1-client-protocol-code.h (source SHA-256: `8a0fc6f86c7bf19e9d5f111e78a0ee3ac5a425f869ee6f9a520462992d4a94f0`)

````text
/*
 * Copyright © 2018 Simon Ser
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/xdg-activation-v1-client-protocol-code.h (source SHA-256: `ddd77f9590376b519c5278dfddb2daf897fd4e6b4a09d9db2c362fb4f25f408a`)

````text
/*
 * Copyright © 2020 Aleix Pol Gonzalez <aleixpol@kde.org>
 * Copyright © 2020 Carlos Garnacho <carlosg@gnome.org>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/fractional-scale-v1-client-protocol-code.h (source SHA-256: `518862beef0ec3dc1a48b53b0dea7eb1762045739adddc4d3c551ec817620e65`)

````text
/*
 * Copyright © 2022 Kenny Levinsen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/deps/mingw/xinput.h (source SHA-256: `4f78d7a9df66252a02342dc874b786ed666f9c32b523ae3734d44d87e09fc0e7`)

````text
/*
 * The Wine project - Xinput Joystick Library
 * Copyright 2008 Andrew Fenn
 *
 * This library is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 2.1 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA 02110-1301, USA
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/fractional-scale-v1-client-protocol.h (source SHA-256: `ea2560e8736086ac6544577a3840d4b3b8f6231e436dbe4f4d2ac7511eae82ee`)

````text
/**
 * @page page_fractional_scale_v1 The fractional_scale_v1 protocol
 * Protocol for requesting fractional surface scales
 *
 * @section page_desc_fractional_scale_v1 Description
 *
 * This protocol allows a compositor to suggest for surfaces to render at
 * fractional scales.
 *
 * A client can submit scaled content by utilizing wp_viewport. This is done by
 * creating a wp_viewport object for the surface and setting the destination
 * rectangle to the surface size before the scale factor is applied.
 *
 * The buffer size is calculated by multiplying the surface size by the
 * intended scale.
 *
 * The wl_surface buffer scale should remain set to 1.
 *
 * If a surface has a surface-local size of 100 px by 50 px and wishes to
 * submit buffers with a scale of 1.5, then a buffer of 150px by 75 px should
 * be used and the wp_viewport destination rectangle should be 100 px by 50 px.
 *
 * For toplevel surfaces, the size is rounded halfway away from zero. The
 * rounding algorithm for subsurface position and size is not defined.
 *
 * @section page_ifaces_fractional_scale_v1 Interfaces
 * - @subpage page_iface_wp_fractional_scale_manager_v1 - fractional surface scale information
 * - @subpage page_iface_wp_fractional_scale_v1 - fractional scale interface to a wl_surface
 * @section page_copyright_fractional_scale_v1 Copyright
 * <pre>
 *
 * Copyright © 2022 Kenny Levinsen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/idle-inhibit-unstable-v1-client-protocol.h (source SHA-256: `7676d4010b406bd610ca8fb8fe4bd718a35b91e0aba56d1fe894906dfeba93d1`)

````text
/**
 * @page page_idle_inhibit_unstable_v1 The idle_inhibit_unstable_v1 protocol
 * @section page_ifaces_idle_inhibit_unstable_v1 Interfaces
 * - @subpage page_iface_zwp_idle_inhibit_manager_v1 - control behavior when display idles
 * - @subpage page_iface_zwp_idle_inhibitor_v1 - context object for inhibiting idle behavior
 * @section page_copyright_idle_inhibit_unstable_v1 Copyright
 * <pre>
 *
 * Copyright © 2015 Samsung Electronics Co., Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/pointer-constraints-unstable-v1-client-protocol.h (source SHA-256: `863eef944b05bf2ac45b64cecd102f28571c3d466e05de9cc941f5d996b2c8a9`)

````text
/**
 * @page page_pointer_constraints_unstable_v1 The pointer_constraints_unstable_v1 protocol
 * protocol for constraining pointer motions
 *
 * @section page_desc_pointer_constraints_unstable_v1 Description
 *
 * This protocol specifies a set of interfaces used for adding constraints to
 * the motion of a pointer. Possible constraints include confining pointer
 * motions to a given region, or locking it to its current position.
 *
 * In order to constrain the pointer, a client must first bind the global
 * interface "wp_pointer_constraints" which, if a compositor supports pointer
 * constraints, is exposed by the registry. Using the bound global object, the
 * client uses the request that corresponds to the type of constraint it wants
 * to make. See wp_pointer_constraints for more details.
 *
 * Warning! The protocol described in this file is experimental and backward
 * incompatible changes may be made. Backward compatible changes may be added
 * together with the corresponding interface version bump. Backward
 * incompatible changes are done by bumping the version number in the protocol
 * and interface names and resetting the interface version. Once the protocol
 * is to be declared stable, the 'z' prefix and the version number in the
 * protocol and interface names are removed and the interface version number is
 * reset.
 *
 * @section page_ifaces_pointer_constraints_unstable_v1 Interfaces
 * - @subpage page_iface_zwp_pointer_constraints_v1 - constrain the movement of a pointer
 * - @subpage page_iface_zwp_locked_pointer_v1 - receive relative pointer motion events
 * - @subpage page_iface_zwp_confined_pointer_v1 - confined pointer object
 * @section page_copyright_pointer_constraints_unstable_v1 Copyright
 * <pre>
 *
 * Copyright © 2014      Jonas Ådahl
 * Copyright © 2015      Red Hat Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/relative-pointer-unstable-v1-client-protocol.h (source SHA-256: `ae911e9aad8b1779ae07bdfc697e1d72f3cc517198f41f9971de5be7c8dea7c6`)

````text
/**
 * @page page_relative_pointer_unstable_v1 The relative_pointer_unstable_v1 protocol
 * protocol for relative pointer motion events
 *
 * @section page_desc_relative_pointer_unstable_v1 Description
 *
 * This protocol specifies a set of interfaces used for making clients able to
 * receive relative pointer events not obstructed by barriers (such as the
 * monitor edge or other pointer barriers).
 *
 * To start receiving relative pointer events, a client must first bind the
 * global interface "wp_relative_pointer_manager" which, if a compositor
 * supports relative pointer motion events, is exposed by the registry. After
 * having created the relative pointer manager proxy object, the client uses
 * it to create the actual relative pointer object using the
 * "get_relative_pointer" request given a wl_pointer. The relative pointer
 * motion events will then, when applicable, be transmitted via the proxy of
 * the newly created relative pointer object. See the documentation of the
 * relative pointer interface for more details.
 *
 * Warning! The protocol described in this file is experimental and backward
 * incompatible changes may be made. Backward compatible changes may be added
 * together with the corresponding interface version bump. Backward
 * incompatible changes are done by bumping the version number in the protocol
 * and interface names and resetting the interface version. Once the protocol
 * is to be declared stable, the 'z' prefix and the version number in the
 * protocol and interface names are removed and the interface version number is
 * reset.
 *
 * @section page_ifaces_relative_pointer_unstable_v1 Interfaces
 * - @subpage page_iface_zwp_relative_pointer_manager_v1 - get relative pointer objects
 * - @subpage page_iface_zwp_relative_pointer_v1 - relative pointer object
 * @section page_copyright_relative_pointer_unstable_v1 Copyright
 * <pre>
 *
 * Copyright © 2014      Jonas Ådahl
 * Copyright © 2015      Red Hat Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/viewporter-client-protocol.h (source SHA-256: `5bafe969401b596ddc3270c7b49b408d8f9bbe60d6a1b64d3dca0419f6b0e397`)

````text
/**
 * @page page_viewporter The viewporter protocol
 * @section page_ifaces_viewporter Interfaces
 * - @subpage page_iface_wp_viewporter - surface cropping and scaling
 * - @subpage page_iface_wp_viewport - crop and scale interface to a wl_surface
 * @section page_copyright_viewporter Copyright
 * <pre>
 *
 * Copyright © 2013-2016 Collabora, Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/wayland-client-protocol.h (source SHA-256: `6adf73ca340109cda2e757a577229b2eef2c01585c6aad296554315ebc50b97d`)

````text
/**
 * @page page_wayland The wayland protocol
 * @section page_ifaces_wayland Interfaces
 * - @subpage page_iface_wl_display - core global object
 * - @subpage page_iface_wl_registry - global registry object
 * - @subpage page_iface_wl_callback - callback object
 * - @subpage page_iface_wl_compositor - the compositor singleton
 * - @subpage page_iface_wl_shm_pool - a shared memory pool
 * - @subpage page_iface_wl_shm - shared memory support
 * - @subpage page_iface_wl_buffer - content for a wl_surface
 * - @subpage page_iface_wl_data_offer - offer to transfer data
 * - @subpage page_iface_wl_data_source - offer to transfer data
 * - @subpage page_iface_wl_data_device - data transfer device
 * - @subpage page_iface_wl_data_device_manager - data transfer interface
 * - @subpage page_iface_wl_shell - create desktop-style surfaces
 * - @subpage page_iface_wl_shell_surface - desktop-style metadata interface
 * - @subpage page_iface_wl_surface - an onscreen surface
 * - @subpage page_iface_wl_seat - group of input devices
 * - @subpage page_iface_wl_pointer - pointer input device
 * - @subpage page_iface_wl_keyboard - keyboard input device
 * - @subpage page_iface_wl_touch - touchscreen input device
 * - @subpage page_iface_wl_output - compositor output region
 * - @subpage page_iface_wl_region - region interface
 * - @subpage page_iface_wl_subcompositor - sub-surface compositing
 * - @subpage page_iface_wl_subsurface - sub-surface interface to a wl_surface
 * @section page_copyright_wayland Copyright
 * <pre>
 *
 * Copyright © 2008-2011 Kristian Høgsberg
 * Copyright © 2010-2011 Intel Corporation
 * Copyright © 2012-2013 Collabora, Ltd.
 *
 * Permission is hereby granted, free of charge, to any person
 * obtaining a copy of this software and associated documentation files
 * (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge,
 * publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the
 * next paragraph) shall be included in all copies or substantial
 * portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT.  IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
 * BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
 * ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/xdg-activation-v1-client-protocol.h (source SHA-256: `159d0005eef9449a818929303ce6563264a9f25a3ec398bd5c601219f02d9d00`)

````text
/**
 * @page page_xdg_activation_v1 The xdg_activation_v1 protocol
 * Protocol for requesting activation of surfaces
 *
 * @section page_desc_xdg_activation_v1 Description
 *
 * The way for a client to pass focus to another toplevel is as follows.
 *
 * The client that intends to activate another toplevel uses the
 * xdg_activation_v1.get_activation_token request to get an activation token.
 * This token is then forwarded to the client, which is supposed to activate
 * one of its surfaces, through a separate band of communication.
 *
 * One established way of doing this is through the XDG_ACTIVATION_TOKEN
 * environment variable of a newly launched child process. The child process
 * should unset the environment variable again right after reading it out in
 * order to avoid propagating it to other child processes.
 *
 * Another established way exists for Applications implementing the D-Bus
 * interface org.freedesktop.Application, which should get their token under
 * activation-token on their platform_data.
 *
 * In general activation tokens may be transferred across clients through
 * means not described in this protocol.
 *
 * The client to be activated will then pass the token
 * it received to the xdg_activation_v1.activate request. The compositor can
 * then use this token to decide how to react to the activation request.
 *
 * The token the activating client gets may be ineffective either already at
 * the time it receives it, for example if it was not focused, for focus
 * stealing prevention. The activating client will have no way to discover
 * the validity of the token, and may still forward it to the to be activated
 * client.
 *
 * The created activation token may optionally get information attached to it
 * that can be used by the compositor to identify the application that we
 * intend to activate. This can for example be used to display a visual hint
 * about what application is being started.
 *
 * Warning! The protocol described in this file is currently in the testing
 * phase. Backward compatible changes may be added together with the
 * corresponding interface version bump. Backward incompatible changes can
 * only be done by creating a new major version of the extension.
 *
 * @section page_ifaces_xdg_activation_v1 Interfaces
 * - @subpage page_iface_xdg_activation_v1 - interface for activating surfaces
 * - @subpage page_iface_xdg_activation_token_v1 - an exported activation handle
 * @section page_copyright_xdg_activation_v1 Copyright
 * <pre>
 *
 * Copyright © 2020 Aleix Pol Gonzalez <aleixpol@kde.org>
 * Copyright © 2020 Carlos Garnacho <carlosg@gnome.org>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/xdg-decoration-unstable-v1-client-protocol.h (source SHA-256: `0d9619c7166d5f19cb46ce407929838a35d3637cab1d87525617f3b8b12de159`)

````text
/**
 * @page page_xdg_decoration_unstable_v1 The xdg_decoration_unstable_v1 protocol
 * @section page_ifaces_xdg_decoration_unstable_v1 Interfaces
 * - @subpage page_iface_zxdg_decoration_manager_v1 - window decoration manager
 * - @subpage page_iface_zxdg_toplevel_decoration_v1 - decoration object for a toplevel surface
 * @section page_copyright_xdg_decoration_unstable_v1 Copyright
 * <pre>
 *
 * Copyright © 2018 Simon Ser
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/xdg-shell-client-protocol.h (source SHA-256: `6f5798f1fbc91fb3d8de3f1b3a8387bf498d1b9568719d927a15b22296a64e4a`)

````text
/**
 * @page page_xdg_shell The xdg_shell protocol
 * @section page_ifaces_xdg_shell Interfaces
 * - @subpage page_iface_xdg_wm_base - create desktop-style surfaces
 * - @subpage page_iface_xdg_positioner - child surface positioner
 * - @subpage page_iface_xdg_surface - desktop user interface surface base interface
 * - @subpage page_iface_xdg_toplevel - toplevel surface
 * - @subpage page_iface_xdg_popup - short-lived, popup surfaces for menus
 * @section page_copyright_xdg_shell Copyright
 * <pre>
 *
 * Copyright © 2008-2013 Kristian Høgsberg
 * Copyright © 2013      Rafael Antognolli
 * Copyright © 2013      Jasper St. Pierre
 * Copyright © 2010-2013 Intel Corporation
 * Copyright © 2015-2017 Samsung Electronics Co., Ltd
 * Copyright © 2015-2017 Red Hat Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a
 * copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice (including the next
 * paragraph) shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
 * THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
 * FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 * DEALINGS IN THE SOFTWARE.
 * </pre>
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/deps/mingw/_mingw_dxhelper.h (source SHA-256: `f3145fbc30022e781a9836d72bd3a628a45495bda5b48f90983e9ae7ebe94b82`)

````text
/**
 * This file has no copyright assigned and is placed in the Public Domain.
 * This file is part of the mingw-w64 runtime package.
 * No warranty is given; refer to the file DISCLAIMER within this package.
 */
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/GLFW/glfw3native.h (source SHA-256: `40035881fdc02a9ab265c2086338f749a0411e1c15b68fada678fa34e8bd8958`)

````text
/*************************************************************************
 * GLFW 3.4 - www.glfw.org
 * A library for OpenGL, window and input
 *------------------------------------------------------------------------
 * Copyright (c) 2002-2006 Marcus Geelnard
 * Copyright (c) 2006-2018 Camilla Löwy <elmindreda@glfw.org>
 *
 * This software is provided 'as-is', without any express or implied
 * warranty. In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would
 *    be appreciated but is not required.
 *
 * 2. Altered source versions must be plainly marked as such, and must not
 *    be misrepresented as being the original software.
 *
 * 3. This notice may not be removed or altered from any source
 *    distribution.
 *
 *************************************************************************/
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/include/GLFW/glfw3.h (source SHA-256: `aa370985f6b493bbe0358a36ab49f5780a6397c0209c6d98a62143dea595b73c`)

````text
/*************************************************************************
 * GLFW 3.4 - www.glfw.org
 * A library for OpenGL, window and input
 *------------------------------------------------------------------------
 * Copyright (c) 2002-2006 Marcus Geelnard
 * Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
 *
 * This software is provided 'as-is', without any express or implied
 * warranty. In no event will the authors be held liable for any damages
 * arising from the use of this software.
 *
 * Permission is granted to anyone to use this software for any purpose,
 * including commercial applications, and to alter it and redistribute it
 * freely, subject to the following restrictions:
 *
 * 1. The origin of this software must not be misrepresented; you must not
 *    claim that you wrote the original software. If you use this software
 *    in a product, an acknowledgment in the product documentation would
 *    be appreciated but is not required.
 *
 * 2. Altered source versions must be plainly marked as such, and must not
 *    be misrepresented as being the original software.
 *
 * 3. This notice may not be removed or altered from any source
 *    distribution.
 *
 *************************************************************************/
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/mappings.h (source SHA-256: `66808cb6f704c7d963a7bae65c02a082f9bb2b5a854d033ce4bd57e22368433f`)

````text
// All gamepad mappings not labeled GLFW are copied from the
// SDL_GameControllerDB project under the following license:
//
// Simple DirectMedia Layer
// Copyright (C) 1997-2013 Sam Lantinga <slouken@libsdl.org>
//
// This software is provided 'as-is', without any express or implied warranty.
// In no event will the authors be held liable for any damages arising from the
// use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not be
//    misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source distribution.
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/x11_window.c (source SHA-256: `7c5a60a507767c5b311437a74e85e547914df631be4f74fa75746063762019da`)

````text
// Decode a Unicode code point from a UTF-8 stream
// Based on cutef8 by Jeff Bezanson (Public Domain)
//
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/init.c (source SHA-256: `61ac2188518cae852bc2bfe0c0f2363dad365f8166099a3aa78a3c9ff8504654`)

````text
// Encode a Unicode code point to a UTF-8 stream
// Based on cutef8 by Jeff Bezanson (Public Domain)
//
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/context.c (source SHA-256: `839a96ff1dcc712129f9bcb9522b8d459b9d6a9b8d029ce32c729c523b9ea012`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2016 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/init.c (source SHA-256: `61ac2188518cae852bc2bfe0c0f2363dad365f8166099a3aa78a3c9ff8504654`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/platform.c (source SHA-256: `117f5e8698eb4f028d7f80ae594fd2b2eb141f55c6122a6c24fd45d1b90209bc`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/platform.h (source SHA-256: `14498e12db44b36f695d5cacd2ca2763fd6d6c4d2c5ef3fdcfa103e79d54b49b`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/vulkan.c (source SHA-256: `5948e9f0b347640473d23f7dfd94cea82375149639f2cbd288491747312d3a33`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2018 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/input.c (source SHA-256: `f5efbc7ec6a38a6c462cd2311dbe51616a4810c4dc8d1a2f735195e5cf72c753`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/internal.h (source SHA-256: `cc6e93eac5b9957c99bc90e95c34d62f184395345de33e554ab77c23a1148ae0`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/monitor.c (source SHA-256: `1671bbe37502419f08b8072e80a13a78494feea6c008ccb2727ebd44a43afd6a`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/window.c (source SHA-256: `c102fd119fd7423ca0cebd17c1fba722d7fe149ef689f09962f6ff79c82dd289`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
// Copyright (c) 2012 Torsten Walluhn <tw@mad-cad.net>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/null_joystick.h (source SHA-256: `ccd7bec0a32632b0fc81d1c085bba3b89e01a9dc4b6b28cde9ea394909fc700d`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/mappings.h (source SHA-256: `66808cb6f704c7d963a7bae65c02a082f9bb2b5a854d033ce4bd57e22368433f`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2006-2018 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
// As mappings.h.in, this file is used by CMake to produce the mappings.h
// header file.  If you are adding a GLFW specific gamepad mapping, this is
// where to put it.
//========================================================================
// As mappings.h, this provides all pre-defined gamepad mappings, including
// all available in SDL_GameControllerDB.  Do not edit this file.  Any gamepad
// mappings not specific to GLFW should be submitted to SDL_GameControllerDB.
// This file can be re-generated from mappings.h.in and the upstream
// gamecontrollerdb.txt with the 'update_mappings' CMake target.
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/null_init.c (source SHA-256: `041b02232d5378e9a3ac3fd12da1e91f8276ef1459303e2094ee820e675a3dd1`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/null_platform.h (source SHA-256: `a93596d5f503d0378de4c347f55d0cd63aad471022a09e3c474f516416520c0e`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2016 Google Inc.
// Copyright (c) 2016-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/null_monitor.c (source SHA-256: `12face8cf2462b5590d8f1ac5502184e0246696f5d5b117dddca3dfbb5658759`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/null_window.c (source SHA-256: `244c6843b93298f97066a5d40324286b9a2de240200ead019470b0e8846595d9`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2016 Google Inc.
// Copyright (c) 2016-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/null_joystick.c (source SHA-256: `5da9465db10034b81b87c3722328bff7dfc3b4a94670cce71672b7d464fd6820`)

````text
//========================================================================
// GLFW 3.4 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2016-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_joystick.h (source SHA-256: `88fb37c0b8bbfea19343e75cfa7090eef7bed67c93bcd789af356ee6c2401210`)

````text
//========================================================================
// GLFW 3.4 Cocoa - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_joystick.m (source SHA-256: `eddd49dc3862a4e16be208c36e540f52fe7f987acf315c4d969eaa6ec32f09b0`)

````text
//========================================================================
// GLFW 3.4 Cocoa - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2009-2019 Camilla Löwy <elmindreda@glfw.org>
// Copyright (c) 2012 Torsten Walluhn <tw@mad-cad.net>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/egl_context.c (source SHA-256: `933150523fa2cf24e8f3e7cd279b0d8308635f769e052f035e4088c4ebc83ccf`)

````text
//========================================================================
// GLFW 3.4 EGL - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/glx_context.c (source SHA-256: `e04e784373d59a3f3b08e3de15c8fbfef61c22b9325a2f3d9922a21cc4dd504e`)

````text
//========================================================================
// GLFW 3.4 GLX - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/linux_joystick.c (source SHA-256: `4e9b95479d3df624d3a5de2b9d708f418cf7e0870c5f05ff2f1a1f88ec84e608`)

````text
//========================================================================
// GLFW 3.4 Linux - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/linux_joystick.h (source SHA-256: `af4e596d0ed1e7eed2699b6567d6ed46157f7c9887bfef15aad5ca15b2a17590`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/xkb_unicode.h (source SHA-256: `b80f53223f42108f4ff0fc51cb8e94c2c9b3ba87888760842958bf43bd4cd19a`)

````text
//========================================================================
// GLFW 3.4 Linux - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2014 Jonas Ådahl <jadahl@gmail.com>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/osmesa_context.c (source SHA-256: `2bb65aa05606a7bb652ee59fb481c45abb223795163fe9c57d808de3e082b09b`)

````text
//========================================================================
// GLFW 3.4 OSMesa - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2016 Google Inc.
// Copyright (c) 2016-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_thread.c (source SHA-256: `898502d86490e8c2e401a3de7a63e7c9566aa99333d676c5c4aaca69b35022b3`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_thread.h (source SHA-256: `10e4679442f9650ec2f423572537a0326e972b8f0c02cfb5fb76c3cd107733cd`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_time.c (source SHA-256: `548431aa62bc3ead2a82c1f049fd53659b7a92c8702abaa9d9176a0583edadb2`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_time.h (source SHA-256: `eccea70e4ada8f6c0fcca9eeb5de01781d2422d28bfe63e795aeb74a7400125c`)

````text
//========================================================================
// GLFW 3.4 POSIX - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_module.c (source SHA-256: `e78e668d01a15695b1d44da660daf8c0fd1696cce1c9b5cca0a2e4e41d67dbe1`)

````text
//========================================================================
// GLFW 3.4 POSIX - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2021 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_poll.c (source SHA-256: `c7fa26ae3a324f9ba15de5bfadfbe356197d4ab667efbf280c41971808e4bce3`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/posix_poll.h (source SHA-256: `4aaeb5ddcc8e8b617463d07d769264d359c8c508a5a9cc9a02b4b49f8dadb3a5`)

````text
//========================================================================
// GLFW 3.4 POSIX - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2022 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/wgl_context.c (source SHA-256: `48c2cad35265f7311146182ec4ad811f551929cb6847ef3a0653abbedec88bf8`)

````text
//========================================================================
// GLFW 3.4 WGL - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/wl_init.c (source SHA-256: `56f0d25a4e665f959749427c7c90498f414c9f515403671027e2da6a0000f77d`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/wl_monitor.c (source SHA-256: `47161326acd583a0787a37927535ffdbe03561923f40a316d9f02b40c28ec16b`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/wl_platform.h (source SHA-256: `fb13b006ad8e13df67a8c267e880324f6f131e1d5cba8b072061e91efa30ee71`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/wl_window.c (source SHA-256: `2f746b505d4e8ac4da3f3ec8f1880bad7f479435a6cd122baad43e33bf672dea`)

````text
//========================================================================
// GLFW 3.4 Wayland - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2014 Jonas Ådahl <jadahl@gmail.com>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_thread.c (source SHA-256: `92ee8a23a6a729c2cff69975d2dfc46d019a2be0494d17b705495cd3042dcd36`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_thread.h (source SHA-256: `cefdd801463eb6901e3ebae02b8204e87e44bd06b2bdc562794dc890886ecc97`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_time.c (source SHA-256: `16c54578192e06205a0dbe675ffc93f73df1d971ace8e81ad2a0ba868c99c4e6`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_time.h (source SHA-256: `dabfeb7478945930cb6c4ded23cfe76f6e6577b68523c5d8ac2acb1b63b35b21`)

````text
//========================================================================
// GLFW 3.4 Win32 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_init.c (source SHA-256: `48d72f04c500d0c0b2e91444429e141972f81ae03072f6d28b26222ac06983c1`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_joystick.c (source SHA-256: `b9c95eed4eeed6eb64864fa031349de83e3421ab7a6c3dcd86cf80fd6b1ad8b1`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_monitor.c (source SHA-256: `4645c2972d8dbf9e11e287b05333f6db585615541217e3dd5d61e9f9fca32e97`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_platform.h (source SHA-256: `f7b39efe79c9731762c9dd9848431e795a66fccc871d7ae72930c5e1aeed7c4c`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_window.c (source SHA-256: `863d05bc961acab390e7d4329ced8a2d36f9b10106b9258376d2f72b0cac297f`)

````text
//========================================================================
// GLFW 3.4 Win32 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_joystick.h (source SHA-256: `41cd05e7477314241c5aca9b840682a4ba4b82c358d401612ab83e5bce167c39`)

````text
//========================================================================
// GLFW 3.4 Win32 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/win32_module.c (source SHA-256: `58f0b4a96584b6ad582e90974b07ecbe9ed8d402cadbe54383afbb0ba37d6a0b`)

````text
//========================================================================
// GLFW 3.4 Win32 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2021 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/xkb_unicode.c (source SHA-256: `57ddc9155229a445c1f6e1542fe3ddce8d4ad0c5d5f7f0c3da476d58fb5f650c`)

````text
//========================================================================
// GLFW 3.4 X11 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2017 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/x11_init.c (source SHA-256: `4c091f59628d3221c6fc2fb3e1451b6147dc86118fa371f06984915581e735e8`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/x11_monitor.c (source SHA-256: `cd1c15c5fc684a58472ded4cbad3509f383c6d6611b4437d8469b54feedbdaab`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/x11_platform.h (source SHA-256: `c35ff50d31c47eee88a7613b8ff47611d6bd892657545b6d44b1911dc7b7466a`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/x11_window.c (source SHA-256: `7c5a60a507767c5b311437a74e85e547914df631be4f74fa75746063762019da`)

````text
//========================================================================
// GLFW 3.4 X11 - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_monitor.m (source SHA-256: `6659a015ce51ca44f0cde5335b8c4e2baacf3a492cccaa883db10a7829a8263f`)

````text
//========================================================================
// GLFW 3.4 macOS - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_time.c (source SHA-256: `a0094a28331172292a7c506935ee01adb0184d18e89e038ea471094781dd0ce4`)

````text
//========================================================================
// GLFW 3.4 macOS - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2009-2016 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_init.m (source SHA-256: `0ce5d37c469a75599257309d3a8d8775d6d1fbc3efa0c99152a77a172069fe1f`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_platform.h (source SHA-256: `b4d608d7fe69fe147ba600aee4700a9e0b40be225ae9abd698f711fdbf944320`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_window.m (source SHA-256: `1c280c56aae81d656f736acb69d3a35b8ac55db234665203df953560dd085a78`)
- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/nsgl_context.m (source SHA-256: `a0d683795265b18cd183be6a185742e44b866ee80c21d7fcc4ad4363c16d2832`)

````text
//========================================================================
// GLFW 3.4 macOS - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2009-2019 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````

### Original source notice

- github.com/go-gl/glfw/v3.4/glfw@v0.1.0-pre.2/glfw/src/cocoa_time.h (source SHA-256: `596db814eaead6d4cf304fa6ae30688362ce98c9ee551eb3871c3cc166876804`)

````text
//========================================================================
// GLFW 3.4 macOS - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2009-2021 Camilla Löwy <elmindreda@glfw.org>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================
````


## github.com/go-logr/logr@v1.4.4

Source: https://pkg.go.dev/github.com/go-logr/logr@v1.4.4

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:tG4xh9yMsRCAiodLVTxyrkzSZ9+o0L1Kg/+cPVcbP/8=`

### github.com/go-logr/logr@v1.4.4/LICENSE

SHA-256: `b40930bbcf80744c86c46a12bc9da056641d722716c378f5659b9e555ef833e1`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "{}"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright {yyyy} {name of copyright owner}

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- github.com/go-logr/logr@v1.4.4/context_slog.go (source SHA-256: `82d65edb4fc6d9c364d663431e706c4fb45922c67b5cb8e6c653edd82234b9ee`)
- github.com/go-logr/logr@v1.4.4/logr.go (source SHA-256: `cce9b9499b3111c29e6cc9a816618b4af3a70fa6d44667134215477191f976bc`)

````text
/*
Copyright 2019 The logr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
````

### Original source notice

- github.com/go-logr/logr@v1.4.4/discard.go (source SHA-256: `720e772ee384154b4b3b4aef0d764332eeb731f8dac7a8e4175adc9a7bf63268`)

````text
/*
Copyright 2020 The logr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
````

### Original source notice

- github.com/go-logr/logr@v1.4.4/funcr/funcr.go (source SHA-256: `03fbeb37f16048ba58539f3a0d48a0422bde28e77dff321a1ce6f77d98978e4c`)

````text
/*
Copyright 2021 The logr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
````

### Original source notice

- github.com/go-logr/logr@v1.4.4/context.go (source SHA-256: `f6cfef95b6e3a080551a40e2c24a9fb4577cbbf0362d42a9ccca1ef33cb19a1f`)
- github.com/go-logr/logr@v1.4.4/funcr/slogsink.go (source SHA-256: `d48dc7009d93743f81f0bf9070b27d9dd822be650dbcaa5e6cefe4dba208cc23`)
- github.com/go-logr/logr@v1.4.4/sloghandler.go (source SHA-256: `4c700ef7443140b8338df5a7463b40954bbeba1b8d1106d3c36abd126050f8f0`)
- github.com/go-logr/logr@v1.4.4/slogr.go (source SHA-256: `b7094c50e8ebe87e6a4241d7efcd6de097520e4bab9a135e06305645c4b758db`)
- github.com/go-logr/logr@v1.4.4/slogsink.go (source SHA-256: `a10143d1b668499a9a38f553849e3bc2e77aef51a9a136991e8d14581ecbef96`)

````text
/*
Copyright 2023 The logr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
````


## github.com/go-logr/stdr@v1.2.2

Source: https://pkg.go.dev/github.com/go-logr/stdr@v1.2.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=`

### github.com/go-logr/stdr@v1.2.2/LICENSE

SHA-256: `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- github.com/go-logr/stdr@v1.2.2/stdr.go (source SHA-256: `da80d2020ce5dc8bc8b8e47c5a5150c85432ec68c39da5b208d6d1c9881da656`)

````text
/*
Copyright 2019 The logr Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
````


## github.com/go-text/render@v0.2.1

Source: https://pkg.go.dev/github.com/go-text/render@v0.2.1

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:qwHhxqGUjjg4L0XyJWj7M7bpY75NZM+kBpv2Yfw5mcg=`

### github.com/go-text/render@v0.2.1/LICENSE

SHA-256: `6b1bf94c77d7faa0ce94f8745595bfaf4c404f8b3873f5c2ce55ca15f8e44eca`

````text
This project is provided under the terms of the UNLICENSE or
the BSD license denoted by the following SPDX identifier:

SPDX-License-Identifier: Unlicense OR BSD-3-Clause

You may use the project under the terms of either license.

Both licenses are reproduced below.

----
The BSD 3 Clause License

Copyright 2021 The go-text authors

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
---



---
The UNLICENSE

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org/>
---
````

### github.com/go-text/render@v0.2.1/testdata/OFL.txt

SHA-256: `02f254589f3cc0e6b751faa152144b632041392b36dce0986dc1abf9d58ac21f`

````text
Copyright 2015-2021 Google LLC. All Rights Reserved.

This Font Software is licensed under the SIL Open Font License, Version 1.1.
This license is copied below, and is also available with a FAQ at:
http://scripts.sil.org/OFL


-----------------------------------------------------------
SIL OPEN FONT LICENSE Version 1.1 - 26 February 2007
-----------------------------------------------------------

PREAMBLE
The goals of the Open Font License (OFL) are to stimulate worldwide
development of collaborative font projects, to support the font creation
efforts of academic and linguistic communities, and to provide a free and
open framework in which fonts may be shared and improved in partnership
with others.

The OFL allows the licensed fonts to be used, studied, modified and
redistributed freely as long as they are not sold by themselves. The
fonts, including any derivative works, can be bundled, embedded, 
redistributed and/or sold with any software provided that any reserved
names are not used by derivative works. The fonts and derivatives,
however, cannot be released under any other type of license. The
requirement for fonts to remain under this license does not apply
to any document created using the fonts or their derivatives.

DEFINITIONS
"Font Software" refers to the set of files released by the Copyright
Holder(s) under this license and clearly marked as such. This may
include source files, build scripts and documentation.

"Reserved Font Name" refers to any names specified as such after the
copyright statement(s).

"Original Version" refers to the collection of Font Software components as
distributed by the Copyright Holder(s).

"Modified Version" refers to any derivative made by adding to, deleting,
or substituting -- in part or in whole -- any of the components of the
Original Version, by changing formats or by porting the Font Software to a
new environment.

"Author" refers to any designer, engineer, programmer, technical
writer or other person who contributed to the Font Software.

PERMISSION & CONDITIONS
Permission is hereby granted, free of charge, to any person obtaining
a copy of the Font Software, to use, study, copy, merge, embed, modify,
redistribute, and sell modified and unmodified copies of the Font
Software, subject to the following conditions:

1) Neither the Font Software nor any of its individual components,
in Original or Modified Versions, may be sold by itself.

2) Original or Modified Versions of the Font Software may be bundled,
redistributed and/or sold with any software, provided that each copy
contains the above copyright notice and this license. These can be
included either as stand-alone text files, human-readable headers or
in the appropriate machine-readable metadata fields within text or
binary files as long as those fields can be easily viewed by the user.

3) No Modified Version of the Font Software may use the Reserved Font
Name(s) unless explicit written permission is granted by the corresponding
Copyright Holder. This restriction only applies to the primary font name as
presented to the users.

4) The name(s) of the Copyright Holder(s) or the Author(s) of the Font
Software shall not be used to promote, endorse or advertise any
Modified Version, except to acknowledge the contribution(s) of the
Copyright Holder(s) and the Author(s) or with their explicit written
permission.

5) The Font Software, modified or unmodified, in part or in whole,
must be distributed entirely under this license, and must not be
distributed under any other license. The requirement for fonts to
remain under this license does not apply to any document created
using the Font Software.

TERMINATION
This license becomes null and void if any of the above conditions are
not met.

DISCLAIMER
THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT
OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL THE
COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL
DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM
OTHER DEALINGS IN THE FONT SOFTWARE.
````


## github.com/go-text/typesetting@v0.3.5

Source: https://pkg.go.dev/github.com/go-text/typesetting@v0.3.5

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:XZPUooClHY0Vf/rFyUyuPRNEkawARaFzLMQcXLSEyPk=`

### github.com/go-text/typesetting@v0.3.5/LICENSE

SHA-256: `6b1bf94c77d7faa0ce94f8745595bfaf4c404f8b3873f5c2ce55ca15f8e44eca`

````text
This project is provided under the terms of the UNLICENSE or
the BSD license denoted by the following SPDX identifier:

SPDX-License-Identifier: Unlicense OR BSD-3-Clause

You may use the project under the terms of either license.

Both licenses are reproduced below.

----
The BSD 3 Clause License

Copyright 2021 The go-text authors

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
---



---
The UNLICENSE

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org/>
---
````

### github.com/go-text/typesetting@v0.3.5/font/testdata/OFL.txt

SHA-256: `b8cb059605ed384a35d449bf50a8f7c0e3cd0b752e5289e6c2e9ce0d2d04ac30`

````text
Copyright 2010-2020 The Amiri Project Authors (https://github.com/alif-type/amiri).

This Font Software is licensed under the SIL Open Font License, Version 1.1.
This license is copied below, and is also available with a FAQ at:
http://scripts.sil.org/OFL


-----------------------------------------------------------
SIL OPEN FONT LICENSE Version 1.1 - 26 February 2007
-----------------------------------------------------------

PREAMBLE
The goals of the Open Font License (OFL) are to stimulate worldwide
development of collaborative font projects, to support the font creation
efforts of academic and linguistic communities, and to provide a free and
open framework in which fonts may be shared and improved in partnership
with others.

The OFL allows the licensed fonts to be used, studied, modified and
redistributed freely as long as they are not sold by themselves. The
fonts, including any derivative works, can be bundled, embedded, 
redistributed and/or sold with any software provided that any reserved
names are not used by derivative works. The fonts and derivatives,
however, cannot be released under any other type of license. The
requirement for fonts to remain under this license does not apply
to any document created using the fonts or their derivatives.

DEFINITIONS
"Font Software" refers to the set of files released by the Copyright
Holder(s) under this license and clearly marked as such. This may
include source files, build scripts and documentation.

"Reserved Font Name" refers to any names specified as such after the
copyright statement(s).

"Original Version" refers to the collection of Font Software components as
distributed by the Copyright Holder(s).

"Modified Version" refers to any derivative made by adding to, deleting,
or substituting -- in part or in whole -- any of the components of the
Original Version, by changing formats or by porting the Font Software to a
new environment.

"Author" refers to any designer, engineer, programmer, technical
writer or other person who contributed to the Font Software.

PERMISSION & CONDITIONS
Permission is hereby granted, free of charge, to any person obtaining
a copy of the Font Software, to use, study, copy, merge, embed, modify,
redistribute, and sell modified and unmodified copies of the Font
Software, subject to the following conditions:

1) Neither the Font Software nor any of its individual components,
in Original or Modified Versions, may be sold by itself.

2) Original or Modified Versions of the Font Software may be bundled,
redistributed and/or sold with any software, provided that each copy
contains the above copyright notice and this license. These can be
included either as stand-alone text files, human-readable headers or
in the appropriate machine-readable metadata fields within text or
binary files as long as those fields can be easily viewed by the user.

3) No Modified Version of the Font Software may use the Reserved Font
Name(s) unless explicit written permission is granted by the corresponding
Copyright Holder. This restriction only applies to the primary font name as
presented to the users.

4) The name(s) of the Copyright Holder(s) or the Author(s) of the Font
Software shall not be used to promote, endorse or advertise any
Modified Version, except to acknowledge the contribution(s) of the
Copyright Holder(s) and the Author(s) or with their explicit written
permission.

5) The Font Software, modified or unmodified, in part or in whole,
must be distributed entirely under this license, and must not be
distributed under any other license. The requirement for fonts to
remain under this license does not apply to any document created
using the Font Software.

TERMINATION
This license becomes null and void if any of the above conditions are
not met.

DISCLAIMER
THE FONT SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO ANY WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT
OF COPYRIGHT, PATENT, TRADEMARK, OR OTHER RIGHT. IN NO EVENT SHALL THE
COPYRIGHT HOLDER BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
INCLUDING ANY GENERAL, SPECIAL, INDIRECT, INCIDENTAL, OR CONSEQUENTIAL
DAMAGES, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
FROM, OUT OF THE USE OR INABILITY TO USE THE FONT SOFTWARE OR FROM
OTHER DEALINGS IN THE FONT SOFTWARE.
````

### github.com/go-text/typesetting@v0.3.5/harfbuzz/LICENSE

SHA-256: `e45d96558ba3d406f95dffed6f22c3be07093475ef617bcfddf88444d6b66788`

````text
MIT License

Copyright (c) 2021 B. KUGLER
Copyright © 2010,2011,2012,2013,2014,2015,2016,2017,2018,2019,2020  Google, Inc.
Copyright © 2018,2019,2020  Ebrahim Byagowi
Copyright © 2019,2020  Facebook, Inc. 
Copyright © 2012  Mozilla Foundation
Copyright © 2011  Codethink Limited
Copyright © 2008,2010  Nokia Corporation and/or its subsidiary(-ies)
Copyright © 2009  Keith Stribley
Copyright © 2009  Martin Hosken and SIL International
Copyright © 2007  Chris Wilson
Copyright © 2006  Behdad Esfahbod
Copyright © 2005  David Turner
Copyright © 2004,2007,2008,2009,2010  Red Hat, Inc.
Copyright © 1998-2004  David Turner and Werner Lemberg

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/buffer.go (source SHA-256: `bbf5ef779355b9677725406405f1824bc12a1011c0fc35935ccf22c9e28f1f31`)

````text
/* ported from harfbuzz/src/hb-buffer.hh and hb-buffer.h
 * Copyright © 1998-2004  David Turner and Werner Lemberg
 * Copyright © 2004,2007,2009,2010  Red Hat, Inc.
 * Copyright © 2011,2012  Google, Inc.
 * Red Hat Author(s): Owen Taylor, Behdad Esfahbod
 * Google Author(s): Behdad Esfahbod */
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/font/cff/parser.go (source SHA-256: `f9b2fcca544c9088520dcf4c343aedcb1ab4c88ce272e81057003cc4552797e3`)

````text
/*Copyright*/
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/fontscan/langset_gen.go (source SHA-256: `c9e40547532e28039386fe5403a8fdad4ef1d3e8986385ba051b4061f2407ee5`)

````text
// Copyright © 2002 Keith Packard
//
// Permission to use, copy, modify, distribute, and sell this software and its
// documentation for any purpose is hereby granted without fee, provided that
// the above copyright notice appear in all copies and that both that
// copyright notice and this permission notice appear in supporting
// documentation, and that the name of the author(s) not be used in
// advertising or publicity pertaining to distribution of the software without
// specific, written prior permission.  The authors make no
// representations about the suitability of this software for any purpose.  It
// is provided "as is" without express or implied warranty.
//
// THE AUTHOR(S) DISCLAIMS ALL WARRANTIES WITH REGARD TO THIS SOFTWARE,
// INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS, IN NO
// EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY SPECIAL, INDIRECT OR
// CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS OF USE,
// DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
// TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR
// PERFORMANCE OF THIS SOFTWARE.
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/font/aat_layout_kern_kerx.go (source SHA-256: `85708fcfd03c848cb18405a883f88e08fd2ea978fe6a7d53298bcdbc1e7d23ba`)
- github.com/go-text/typesetting@v0.3.5/font/aat_layout_morx.go (source SHA-256: `0bcbe4dfb66a54c56bb9120e108fdca60bc877b83df0a365f2a6b0f1639e01a4`)
- github.com/go-text/typesetting@v0.3.5/font/bitmaps.go (source SHA-256: `51cb4b6d22698403ee1672df8a88e3dc331f02eed15b246b899b165a2edee3d5`)
- github.com/go-text/typesetting@v0.3.5/font/cff/cff_gen.go (source SHA-256: `1ad9247322f7b331425ab9ddc24ea3bab2a679f1bf4a86d797c40a8c233fc9ca`)
- github.com/go-text/typesetting@v0.3.5/font/cff/charsets.go (source SHA-256: `63216673811b34c5964129a18d70645692974a330b72a527afd1fb59cdfda416`)
- github.com/go-text/typesetting@v0.3.5/font/cff/charstring.go (source SHA-256: `eaa58157a896d2ad2d0838eff2ab353bdd66bc3da7585353aa62d941c44379e7`)
- github.com/go-text/typesetting@v0.3.5/font/cff/interpreter/charstrings.go (source SHA-256: `b193c16d4e2d7db6fc21929732f0b0578f9a9e29e09b6223633335d345463916`)
- github.com/go-text/typesetting@v0.3.5/font/cff/interpreter/interpreter.go (source SHA-256: `2657bb237985813ad84f339312aa92e251073567b7272936be74e98342c570be`)
- github.com/go-text/typesetting@v0.3.5/font/cff/parser.go (source SHA-256: `f9b2fcca544c9088520dcf4c343aedcb1ab4c88ce272e81057003cc4552797e3`)
- github.com/go-text/typesetting@v0.3.5/font/cmap.go (source SHA-256: `42be319073560b5a69606ca6867ce347bb1e2082af734c418ecbd676105b0b9c`)
- github.com/go-text/typesetting@v0.3.5/font/cmap_arabic_pua_table.go (source SHA-256: `145575dc9447f879a851e8d759024fd8f61cc0f6684b32aa07b83bf9040fe4f6`)
- github.com/go-text/typesetting@v0.3.5/font/font.go (source SHA-256: `1495a51931e802a8bd187e717d63cca0fe3aae3afdeb76d1b57958b7c22ca6f8`)
- github.com/go-text/typesetting@v0.3.5/font/glyphs.go (source SHA-256: `6bb9d1f3eaafe4146e02828abe940d74f8298ed6d1afd735a5346e4ba75ef5be`)
- github.com/go-text/typesetting@v0.3.5/font/metrics.go (source SHA-256: `8aa6889621b07929eaed03a96d6903c6ee1aad2fbc2314318300a63592ea27a0`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/opentype.go (source SHA-256: `18230d25c20659c6702252cdee044277e98f5327f8b9c75a6a10a95b390966d5`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/reader.go (source SHA-256: `6da2b334e7355de6ac193aed1ab2e4a07b7f0f52404532ac2ba9745fc45e5f60`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/reader_otf.go (source SHA-256: `1241f4140b4dcd253b1c5556ff0761007564809db9f9b23cb418de30ed47defb`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/reader_woff.go (source SHA-256: `3220808624dde972fb5da2522a1f2a6be8ef1e2c05e29df8c318e100628fbac8`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_ankr_gen.go (source SHA-256: `e5f5db218114a2c25e6a219d43843aef432636465ccd1542e77717712324dc93`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_ankr_src.go (source SHA-256: `c0be963bf3ad11b9a9a348db9630c582d51ab1263225bd14daee21ab7509267d`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_common.go (source SHA-256: `f3afea2fd2bdf4c417cd491216520fd5d0dca1fb27267d85e6ff723bcceff4d3`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_feat_gen.go (source SHA-256: `936e077587823fe5b4bd3dc50cb2c4269050f140b0d7c133933c59e7db716776`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_feat_src.go (source SHA-256: `f1bd23ea70d388219568e1819cc23058847222a9d882daabfc5b59571fcce241`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_kerx_gen.go (source SHA-256: `a6335e6d77c6da72a313abe92de5a3fe844db25b5b3dd1cd820622a2a547130d`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_kerx_src.go (source SHA-256: `a14b95d829bcafca539ab84272530d3a4bb319f430019eb8fe60e66b06b0d51c`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_ltag_gen.go (source SHA-256: `1503c663c9c0c100d0fac6c2ea93c53344d95b7b7bd78ca40c1d3a8badeab8e2`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_ltag_src.go (source SHA-256: `6e6acf3e27cb2142050a22848eb1b38d254a57b52752ba8342eca6104204606f`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_mortx_gen.go (source SHA-256: `98ea82fac58d66435d369afe39d0e2ecd0493d3ab83d3b6841f1da66327cac82`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_mortx_src.go (source SHA-256: `4df290bc460fb5366f6d45f69f81613856d90c615ef9d7a51fa67cfd4abc17d4`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_properties.go (source SHA-256: `216088825ee4171f9e06011d12f4cfc48e2c927beb523450222ef2980b68b249`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_trak_gen.go (source SHA-256: `1498d4c9355b35c2c21b15b340e33991636e70c952fd5281e1740b1e5c2d34cf`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/aat_trak_src.go (source SHA-256: `ce26170b15acf40776c4e3a659c9402139bd01e325f87e53dd427c0be559ae4f`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/cmap_gen.go (source SHA-256: `2c5a4f93aeb7c2b7db2084eb3510ef59b212cb56e457bcf0485af13187aa8cc2`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/cmap_src.go (source SHA-256: `9ec359435012ac3d763cc09a28e387e2f00ebf9e40e6c7963d2fdb8fb450709d`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs.go (source SHA-256: `fe90e5cc9ba3de8904b4cddfbf37f145f87324b4ac11584937db34c3591b7f5d`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_bitmap_gen.go (source SHA-256: `cfa7e5b965cbda2640a7fd91938bb6602fea50b1ae3ca95d34d28ce30ec2e35d`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_bitmap_src.go (source SHA-256: `d87d8e73fee7f8c0b9e2a95f72d1c21b17c02e11cdd3aba5989aaa4f36f73402`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_colr_gen.go (source SHA-256: `8b26f88d2fed31aa82793af572e37e1c7564c720f24129a6e8f654fb833f1a97`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_cpal_gen.go (source SHA-256: `47c023320f16fbac83e4303707e667c71742e0c9125a14b081daa82cce3bc569`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_glyf_gen.go (source SHA-256: `9161159a4a87bed9561b6aa15ae3db3eb46b282d99a1d2bcb21e9f3b99c14f80`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_glyf_src.go (source SHA-256: `4489dbda11517de418a7ba3f21c598abb390c8dfeeb242c1272bd2cf0c522ef1`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_misc_gen.go (source SHA-256: `d520a349a2b41648aee234566803f3cbcc450636ac51ed7c6b462c01791aaf3d`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_misc_src.go (source SHA-256: `2afcd55c6d194683e9201c0289d3fe7bab0c1ce1578f97598b975012eed87750`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_sbix_gen.go (source SHA-256: `9ac950f591c8ce4bd851acdbb0669bc10068ddfa92f30dc28749aad4a6a4a529`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/glyphs_sbix_src.go (source SHA-256: `140c2f3d6c27d94eb9d69c2ec475813a2a32a1ba5eba19fe8bd4595497b5e4e6`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/head_gen.go (source SHA-256: `0ec0082ebb3364f82e42a3487934a1ef31e37982bd2b2af1952a942e157b9cf3`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/head_src.go (source SHA-256: `530aa8dae5a3f9e30857a7e30b3b2d597e5b6c08980db3987d67e81cf10a1dad`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/hhea_vhea_gen.go (source SHA-256: `1e0aead6745df2362504c6518a38f9e586dfea29ccf66d0fbc8031b3a99b8486`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/hhea_vhea_src.go (source SHA-256: `ceaf0286846326fdb0dd62278a40c55b560717d3bf0bb5448856a1f1b2570baa`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/hmtx_vmtx_gen.go (source SHA-256: `1f8643ca0e38a61309f7ea05b60b1b2b6aa0c7d4732d1d756ec1db6197c75851`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/hmtx_vmtx_src.go (source SHA-256: `66da9bdeb0c3c14ee35414fb659f512f5abb69a019849d35c9fb160484ae7a70`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/kern.go (source SHA-256: `b700c89ea6354e915859115c41ef39b39d93bab9a2eb87982ccaeecb60c490cd`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/kern_gen.go (source SHA-256: `ae00fd5443a7b5520f4cec3ec3bbe6ac7aed29130bf8c20ee95dda4e2c353df0`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/kern_src.go (source SHA-256: `5da1698ad793a4004740e4b0e7626846c1b45818a750467318a1ded2ff347949`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/maxp_gen.go (source SHA-256: `5eb0a842c451865bb8ddfddcc3ac5cd055e92c35aaa4db3096fc5135c9f7bbaf`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/maxp_src.go (source SHA-256: `f6d1e6cbc63f875c97a8b3eacd3271557dd492ae3fc926a9b9a602127c8d9e5a`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/name_gen.go (source SHA-256: `42ae50ff9a37a782c198599169e633f8a0df32068871651bc285d7ec67dc1817`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/name_src.go (source SHA-256: `43c2223255874040488d0c2f1f313f6f1e1b89070f996c915a4c7bca3cda2859`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/os2_gen.go (source SHA-256: `36c4ef42af344d087ba87f6cec8c23f92de9ebe5c142f8015528192844227c46`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/os2_src.go (source SHA-256: `436901cdbeccbfd44a99778bb7959a383dc2e377df2b4d149b5e0656d7f111e6`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_gdef_gen.go (source SHA-256: `3c931e096abdf5c968080ac305fad80ec946fb819e350e193b14bf1d3d5e279a`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_gdef_src.go (source SHA-256: `462efbeccc6954a251c9051fd64f673c6f894ea897988c38120971e4ceb65139`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_gpos_gen.go (source SHA-256: `88465ac020f426e65e49095c276778916ef48e348807ed9e2a1f8664fa051db2`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_gpos_src.go (source SHA-256: `cb00a335cbcdda6710cc65baeb37fea84672a4ddc5608d358c2784d1469c561b`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_gsub_gen.go (source SHA-256: `65f764b0d60190ad497e7ac3ca11ead5407c1df9b29b5d9f104cd152be88a709`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_gsub_src.go (source SHA-256: `a01b55cb5f1e2775fc9260a4f8d380aae901ddb4e9d1a2c0af5765c7b5118a7e`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_layout.go (source SHA-256: `10089703f5ad4d159d4b3c9b63ce954a8479fc9c09a73d644b51eb117123f74a`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_layout_gen.go (source SHA-256: `4c2b6a72e3a45ab652c9bf8613b0e851ba50fc30c654791f34107ab5c1eb0778`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_layout_src.go (source SHA-256: `02e4293eeff8cb1917294a5812d22f12c7fb1f217a6a14ea9977a68a2fa949e9`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/ot_properties.go (source SHA-256: `a87250822d46d015d8ec5562f1e066caf2ff7f69b79135e779737cd9f941b6f8`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/post_gen.go (source SHA-256: `8041e465a3e7fc159b4a9fbae7fa9ac062cf581b52b4e1d3dd441a671ec467fa`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/post_src.go (source SHA-256: `da6f09036d1ec2c48ad9701daafac99d7c618cd5c74efe1eab466f7043f0a8aa`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/tables.go (source SHA-256: `cd5529ef90e04c68c7a062aaa5536ffb4d71b158d790cdcc8825cd709d644bed`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/xvar_gen.go (source SHA-256: `b1bcf1a69185a8fb0651f1025597d8addd96554d4303e7db9f84990e3786def1`)
- github.com/go-text/typesetting@v0.3.5/font/opentype/tables/xvar_src.go (source SHA-256: `a38e64792c256dd22581453592ef25779cb118bed6018384f9c743546c1fdf36`)
- github.com/go-text/typesetting@v0.3.5/font/os2.go (source SHA-256: `c2ec49ed025dd494eadb5a72ab443d004c3f0c0b74bd8018f4fca20aefb6e2bc`)
- github.com/go-text/typesetting@v0.3.5/font/ot_layout.go (source SHA-256: `c8cc862cfd47f3db3a4afdb19f68511ebe5c73f5e5609c1b8709d8b291faf8fe`)
- github.com/go-text/typesetting@v0.3.5/font/post.go (source SHA-256: `8541e9d852b7bff80237e3772efcc2c03bdf3d6e5fa9ed0dd8f2c6f67680835a`)
- github.com/go-text/typesetting@v0.3.5/font/renderer.go (source SHA-256: `920f3ba6f8beab74fe5e709364ce27614b77a6b39d9b79f11082f5393fd731cf`)
- github.com/go-text/typesetting@v0.3.5/font/svg.go (source SHA-256: `52f9b3dda1a1389b41d07a9a635029a92422ba6cca0a2b15da34934923f76d83`)
- github.com/go-text/typesetting@v0.3.5/font/variations.go (source SHA-256: `110fc702d98407602a170f4609cb420f8b510f2721c5df79f7306452e7abccd3`)
- github.com/go-text/typesetting@v0.3.5/fontscan/fontconfig.go (source SHA-256: `aeb3e95b1d4240f75c1b9451bb0f9fb865842c149e2f2ed1d7d0d4a52721a1d9`)
- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_arabic_table.go (source SHA-256: `949f14773955c3995e39ef5a2778c6f1ddc1bd86af41a723e7b904253df39bff`)
- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_indic_table.go (source SHA-256: `58ba82ea50f8bc5c25d7df81b697b6de3a21df489bfa94242880fe28d822891e`)
- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_use_table.go (source SHA-256: `0f528f3da8a801de8c1a159b45456ef1c556b340b3bc2ef4df03c9a40bcb8bb9`)
- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_vowels_constraints.go (source SHA-256: `91a6730c70bfd698b023b7919e618bbdac20ae443d4d44c82930eb13985ce059`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/bidi_class.go (source SHA-256: `93cd47740a4828f00dc92eff230dd554790dc9b652930e3b590b42a43df332ce`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/combining_classes.go (source SHA-256: `c443db6378eae91df80fc5b017f8e5a51e61f5851f2ed5366d24c2f7d12a4c93`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/decomposition.go (source SHA-256: `a06c6d6136dbd8f0fe9efb469d481ed81248467d18a90ef231e133bd26e4efc0`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/east_asian_width.go (source SHA-256: `619cb07417d0a4ff878fc949060796de66bc84d91b1c15724473262150931c01`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/emojis.go (source SHA-256: `ee6e5a85e46f70d03a1492467ff8148972a9a28e906aa8a2c241e4f6d374d8ba`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/general_category.go (source SHA-256: `02d3ba4107863a52bf07c309159194c48b3fa3b58c2e429e113d0c8ab8799cc6`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/grapheme_break.go (source SHA-256: `cb0daae6cc5e55e557cc973fe0e77b6fc48d1e92d0662099bb602d80486f08a9`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/indic_conjunct_break.go (source SHA-256: `30a9ba1d9233a7a906b80177addf128ab2e99c336d617bad48b832e4eedf829f`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/line_break.go (source SHA-256: `d33c7dda52dedc52e6ab16dd0f22fcb99cf5335475a02e006ffc0399e0324c43`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/mirroring.go (source SHA-256: `20652652fb085857968d82d1a19c9906a7c6de82f4da6d7ee50fc25194ea074c`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/vertical_orientation.go (source SHA-256: `497a892e33004cd8bdc4e4343ccefaee717e5b691bdd2905f6c671e071d5804e`)
- github.com/go-text/typesetting@v0.3.5/internal/unicodedata/word_break.go (source SHA-256: `df74e485d50e5a04ac3c7b71a517dac710698ae0c5a9c0fc1a6efeb270bae9d9`)
- github.com/go-text/typesetting@v0.3.5/language/language.go (source SHA-256: `375fe83a1404561f1fc16c4a6a8d0af0571ad8cb802b688f4fbb34e5903bfa97`)
- github.com/go-text/typesetting@v0.3.5/language/scripts_table.go (source SHA-256: `02a573fca71f5f55184e69036e90fab42e198bdf39c5541f44aec2c1846e1b5e`)
- github.com/go-text/typesetting@v0.3.5/segmenter/segmenter.go (source SHA-256: `9efb82fffe1d41866251889944c0316438a473f9a962c2d440a20322934e180f`)
- github.com/go-text/typesetting@v0.3.5/segmenter/unicode14_rules.go (source SHA-256: `12dc7ba3cf6da713f89703a46786477a5d032aab31539bd1e85739a36c118148`)
- github.com/go-text/typesetting@v0.3.5/segmenter/unicode29_rules.go (source SHA-256: `413f880fd2749deb7c2788c9a5e5ad35f1f0ce53c1b378b7407a4609acdd34b6`)
- github.com/go-text/typesetting@v0.3.5/shaping/input.go (source SHA-256: `349ea57a14b8b704ca1dbdd6288c4cb3a7e24ca1a44c17ec338eeeccf409f4ad`)
- github.com/go-text/typesetting@v0.3.5/shaping/output.go (source SHA-256: `f99dc329c5bbf190255c1f32cb844703529ac885e1763f7beba5bfce3f6f820b`)
- github.com/go-text/typesetting@v0.3.5/shaping/paired_delims_table.go (source SHA-256: `e32c2617811099c5c7eebfea627ac77249e4d5aca6bc2d246a5189ac1222cd72`)
- github.com/go-text/typesetting@v0.3.5/shaping/shaping.go (source SHA-256: `79cbd39df8c42d643406fcf2f68f700b5702faa72667de49787c0364fcfc0c61`)

````text
// SPDX-License-Identifier: Unlicense OR BSD-3-Clause
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape.go (source SHA-256: `cf2985246e6d010f15fd19495256f30e8abc02f54f29dd6fbddde3585fabc532`)

````text
// Support functions for OpenType shaping related queries.
// ported from src/hb-ot-shape.cc Copyright © 2009,2010  Red Hat, Inc. 2010,2011,2012  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/aat_layout.go (source SHA-256: `42603d75b321b296ed0fe844aaf5ec6a4d0840181fd67f50012060bdfc776a55`)

````text
// ported from harfbuzz/src/hb-aat-layout.h  Copyright © 2018 Ebrahim Byagowi, Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/aat_map.go (source SHA-256: `fdb813df5457558ee1aaf0df88b7954c45223b56bed36075de38e23a43b495de`)

````text
// ported from harfbuzz/src/hb-aat-map.cc, hb-att-map.hh Copyright © 2018  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_layout_gpos.go (source SHA-256: `2f620e110f94110bd2d04477e12af7a854567259461e7beb530ac1512719d977`)

````text
// ported from harfbuzz/src/hb-ot-layout-gpos-table.hh Copyright © 2007,2008,2009,2010  Red Hat, Inc.; 2010,2012,2013  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_layout_gsubgpos.go (source SHA-256: `f1a6ad49aa29c0973ded23e6738da785f65a56113511d51450efbfdd2f2fdae8`)

````text
// ported from harfbuzz/src/hb-ot-layout-gsubgpos.hh Copyright © 2007,2008,2009,2010  Red Hat, Inc. 2010,2012  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_map.go (source SHA-256: `73620865da79797407f5cf185d03abb4ba81015e1aed9ec210e4d8d2442e4712`)

````text
// ported from harfbuzz/src/hb-ot-map.cc, hb-ot-map.hh Copyright © 2009,2010  Red Hat, Inc. 2010,2011,2013  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_arabic_win1256.go (source SHA-256: `f84a2ea55e2a4657f74c5aa96d0705da79d918e9a0baedec2fab0b8006236cec`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-arabic-win1256.hh Copyright © 2014  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_arabic.go (source SHA-256: `8576d78d318ae50ba0370b1328c6cf9b72c92f9f2a1561613b6a884ee0101fdd`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-arabic.cc, hb-ot-shape-complex-arabic-fallback.hh Copyright © 2010,2012  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_hangul.go (source SHA-256: `2f2563e4bdc7897b3ca18ae42fcc3193dca048d79e4b4f92cf5c77e2b11f4cbe`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-hangul.cc Copyright © 2013  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_hebrew.go (source SHA-256: `ceae3cd905c2c86fa45dde2711d6e1deb5ba81f354ba85dedd59445f46113d8d`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-hebrew.cc Copyright © 2010,2012  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_indic_machine.go (source SHA-256: `8cd5c6e4286538cb8b4eacb6eb8f6ad060f0fa1c6ffb1df9893a3b5ea3c326d7`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-indic-machine.rl Copyright © 2015 Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_indic.go (source SHA-256: `ef844541e4f38ed61f4ebc83af94e598f52273332ba5cd6f2b8a0aa48a5c2fca`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-indic.cc, .hh Copyright © 2011,2012  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_khmer_machine.go (source SHA-256: `29732d268eb1cabe141e222dfcea10be6dceba7c6cdb4e413fde40d151cf6ff0`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-khmer-machine.rl Copyright © 2015 Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_khmer.go (source SHA-256: `6d99b850135cc95e54169294a8a5ea10b2115d5b3da01aac30beef6806f104aa`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-khmer.cc Copyright © 2011,2012  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_myanmar_machine.go (source SHA-256: `42b77acc167b99ef5950990940d3105f358668e28478132a820eff9f6481f650`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-myanmar-machine.rl Copyright © 2015 Mozilla Foundation. Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_myanmar.go (source SHA-256: `013eb9d031ac5d65b4ccb276daffb67468ff83905dcf8bee479eb8ff0af18543`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-myanmar.cc, .hh Copyright © 2011,2012,2013  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_thai.go (source SHA-256: `2bc6a52fde8d073e232027bdb4bf84e3ff4690d1b9b9dee855f0f2790ece70b7`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-thai.cc Copyright © 2010,2012  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_use_machine.go (source SHA-256: `03f27c97565918900214277831457442933ad6108305550399683028bdb18e17`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-use-machine.rl Copyright © 2015 Mozilla Foundation. Google, Inc. Jonathan Kew Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_use.go (source SHA-256: `41210623f838dfc7a05804fa8741c77d3e85f7fb9ad2dbabff97fe65ef15b577`)

````text
// ported from harfbuzz/src/hb-ot-shape-complex-use.cc Copyright © 2015  Mozilla Foundation. Google, Inc. Jonathan Kew, Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_fallback.go (source SHA-256: `229a8858fadb836d74d7b27d54ed57d75953992a25ef289ba0f1bd3438390f16`)

````text
// ported from harfbuzz/src/hb-ot-shape-fallback.cc Copyright © 2011,2012 Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_shape_normalize.go (source SHA-256: `e868c3c6dc7b94b7cd1f80bf3749992ebf894de1240f0f29fbcdbb08401f6627`)

````text
// ported from harfbuzz/src/hb-ot-shape-normalize.cc Copyright © 2011,2012  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_tag.go (source SHA-256: `24c59a53d52657bc44c46fef024b0dcd2aa8f547ade8cce05b7a9d5017df4ebd`)

````text
// ported from harfbuzz/src/hb-ot-tag.cc Copyright © 2009  Red Hat, Inc. 2011  Google, Inc. Behdad Esfahbod, Roozbeh Pournader
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/shape.go (source SHA-256: `7786182e22b731bf5e570ade159a17525c1d57d5e444f9214a7653b50e51fb34`)

````text
// ported from harfbuzz/src/hb-shape.cc, harfbuzz/src/hb-shape-plan.cc Copyright © 2009, 2012 Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/font.go (source SHA-256: `c9166785892a409ff421a2c92180da4aa2b3c627240c5e63ef36c6b4de003a7b`)

````text
// ported from src/hb-font.hh, src/hb-font.cc  Copyright © 2009  Red Hat, Inc., 2012  Google, Inc.  Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/ot_layout.go (source SHA-256: `3281a6c1b3771471b49d30439f95e3aef6a1a5ba6fe1489b214de5a16c03d40e`)

````text
// ported from src/hb-ot-layout.cc, hb-ot-layout.hh
// Copyright © 1998-2004  David Turner and Werner Lemberg
// Copyright © 2006  2007,2008,2009  Red Hat, Inc. 2012,2013  Google, Inc. Behdad Esfahbod
````

### Original source notice

- github.com/go-text/typesetting@v0.3.5/harfbuzz/set_digest.go (source SHA-256: `c558844607891244b47c61ea93affc08fbfdddccd48acb63767531616c9682d9`)

````text
// ported from src/hb-set-digest.hh Copyright © 2012  Google, Inc. Behdad Esfahbod
````


## github.com/godbus/dbus/v5@v5.2.2

Source: https://pkg.go.dev/github.com/godbus/dbus/v5@v5.2.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:TUR3TgtSVDmjiXOgAAyaZbYmIeP3DPkld3jgKGV8mXQ=`

### github.com/godbus/dbus/v5@v5.2.2/LICENSE

SHA-256: `e4646a82a976369d7ae8f6ed5c11d35dc0af18433a8ccc24c85b459ad8b95128`

````text
Copyright (c) 2013, Georg Reinke (<guelfey at gmail dot com>), Google
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions
are met:

1. Redistributions of source code must retain the above copyright notice,
this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright
notice, this list of conditions and the following disclaimer in the
documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/google/s2a-go@v0.1.11

Source: https://pkg.go.dev/github.com/google/s2a-go@v0.1.11

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:7v1V1TraNGP63WptsAkNqr2w8jgm4vM6A6qx6Lsw9iM=`

### github.com/google/s2a-go@v0.1.11/LICENSE.md

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- github.com/google/s2a-go@v0.1.11/internal/authinfo/authinfo.go (source SHA-256: `04a18d9bd992082525496500dc308243fbd2109d8805c4171e32b9ff800c4a6f`)
- github.com/google/s2a-go@v0.1.11/internal/handshaker/handshaker.go (source SHA-256: `884bb8fe9ca1e644a56b05aef4a18d29a4ddcb085a5280b18f25d741e7ddea17`)
- github.com/google/s2a-go@v0.1.11/internal/handshaker/service/service.go (source SHA-256: `57550f0e6903b2c601d806b775fa0abbd157c2b2f52b260eaf51614113f09c14`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/aeadcrypter/aeadcrypter.go (source SHA-256: `20a929bb1ee91d10aa5fa2acd4d7a3af6866e880740ec1e22c02bf2b6d89650f`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/aeadcrypter/aesgcm.go (source SHA-256: `b84e6f7e105194e022d02b72005fcac11d30efdf31ee560d1d7cb7f0340b6f1d`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/aeadcrypter/chachapoly.go (source SHA-256: `9db845bcda6065b4120fde7b5c331d3f674934d1b7837249727fdd818dfe99a2`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/aeadcrypter/common.go (source SHA-256: `7f89197b10aeb59a5653c2dcfee84a1269266befeb9377ddb6290a5062277ae5`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/halfconn/ciphersuite.go (source SHA-256: `8e95a16af89152a81c18859ca73f7c5f2eeb5c474c11890ad9d46116c5cd600a`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/halfconn/counter.go (source SHA-256: `5ebf8d3746182f5b138c2107185ac87347760f3bc4b3a34590caef5ee9a7a985`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/halfconn/expander.go (source SHA-256: `f668317cc2c625ee04c042dc16e768df29d0ba84d2913c36ce7b991487e610b6`)
- github.com/google/s2a-go@v0.1.11/internal/record/internal/halfconn/halfconn.go (source SHA-256: `39cc88c36500f36df816a0d9b0c0f73bd500e7a68210892d4d0914dc63952440`)
- github.com/google/s2a-go@v0.1.11/internal/record/record.go (source SHA-256: `f04da677efedaa07b925b70bcaadd90359b63515bff48f410b782edeaef3a419`)
- github.com/google/s2a-go@v0.1.11/internal/record/ticketsender.go (source SHA-256: `c5bb4a6419935bd216ce21fdb3eed6fe08abb012e7bd5cc65a50c73a6d62e33b`)
- github.com/google/s2a-go@v0.1.11/internal/tokenmanager/tokenmanager.go (source SHA-256: `2bfabb445a3dedffdbb210bc6d13b56411f0df0d8d5dfa3210fd61151cc5dfbf`)
- github.com/google/s2a-go@v0.1.11/s2a.go (source SHA-256: `ca5ceaa0acbafa589d233cc62a008ee6593d3291ca23dd4678a407e8c1b6dcea`)
- github.com/google/s2a-go@v0.1.11/s2a_options.go (source SHA-256: `beac15a15abceca8d96ceea1aa36de01e49930c33c22d0f75907489b8cbde8fd`)
- github.com/google/s2a-go@v0.1.11/s2a_utils.go (source SHA-256: `5dc66b722a1732c83caa0ab53f8e95ec4d5b368eaaa2b879b571ca720b8853e7`)

````text
/*
 *
 * Copyright 2021 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- github.com/google/s2a-go@v0.1.11/internal/v2/certverifier/certverifier.go (source SHA-256: `8a36a2dfc8cd73731f34a872d8fd64f4b51d49d8bff967bc3160ff71ce988cb9`)
- github.com/google/s2a-go@v0.1.11/internal/v2/remotesigner/remotesigner.go (source SHA-256: `2ba7760a475eda6fb753184b93fd6e4c3dc9727496ac0bf3b4d26ed589d48686`)
- github.com/google/s2a-go@v0.1.11/internal/v2/s2av2.go (source SHA-256: `34713eb675c4bd50bec4b083aca747de42e8bce18fcbe1f5f37aefdcc31fea56`)
- github.com/google/s2a-go@v0.1.11/internal/v2/tlsconfigstore/tlsconfigstore.go (source SHA-256: `394e12878fdede499362c10b79c3abc4400c3aa6f0efbf92b6fea09c503e3858`)

````text
/*
 *
 * Copyright 2022 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- github.com/google/s2a-go@v0.1.11/fallback/s2a_fallback.go (source SHA-256: `e65e77f5a208781f5ea48ddefad12e1961f0bed55c026a303de3d7b3e39638c3`)
- github.com/google/s2a-go@v0.1.11/retry/retry.go (source SHA-256: `533ade116f984f51a6a8056a4a5abf23b4f16e1df2be3e297a43c33cf3bc8d7a`)
- github.com/google/s2a-go@v0.1.11/stream/s2a_stream.go (source SHA-256: `29cb83ea1f0a504ff5b362a4142f99e8725a1dbf1296ee70a9e654c8ad80272e`)

````text
/*
 *
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- github.com/google/s2a-go@v0.1.11/internal/proto/common_go_proto/common.pb.go (source SHA-256: `9eaf8ee61d5c1de046f73e1fb70d43003e843024de9458410a1c59287399bfc9`)
- github.com/google/s2a-go@v0.1.11/internal/proto/s2a_context_go_proto/s2a_context.pb.go (source SHA-256: `6d746c15d0e7a134309c65b8c739eaa4bd01212f67374c4a2023ea7c2f2502cd`)
- github.com/google/s2a-go@v0.1.11/internal/proto/s2a_go_proto/s2a.pb.go (source SHA-256: `575150ba4d0296022d5fa633908cbcc4cf505887da72d54dea655fa2123f371b`)
- github.com/google/s2a-go@v0.1.11/internal/proto/s2a_go_proto/s2a_grpc.pb.go (source SHA-256: `e5fb93da089a24dfec9eb0ba6ce8ba546ac52c43460726774ef9df85d40a5f72`)

````text
// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- github.com/google/s2a-go@v0.1.11/internal/proto/v2/common_go_proto/common.pb.go (source SHA-256: `a7e0683888310f89b5446e1ae5ac66a09c7fd3762fd535d015ac7307d06936a0`)
- github.com/google/s2a-go@v0.1.11/internal/proto/v2/s2a_context_go_proto/s2a_context.pb.go (source SHA-256: `e0c6c00c01270e7fd45ba47c18eb04b992adce8f5d45e968df81fca1fa4038cc`)
- github.com/google/s2a-go@v0.1.11/internal/proto/v2/s2a_go_proto/s2a.pb.go (source SHA-256: `3f9324a6586e51ff6957b3a53db5c0d168cb7a9520b9196248310b59818888ac`)
- github.com/google/s2a-go@v0.1.11/internal/proto/v2/s2a_go_proto/s2a_grpc.pb.go (source SHA-256: `fd9edb5a87cb0f29f42ebdefe7af46ffd540b3419c3791a92cc34fee05d452c0`)

````text
// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## github.com/google/uuid@v1.6.0

Source: https://pkg.go.dev/github.com/google/uuid@v1.6.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=`

### github.com/google/uuid@v1.6.0/LICENSE

SHA-256: `0a8d61ed3cbfd5312326e8126c31ce9c627a283adc99131b56896d29ada04b2d`

````text
Copyright (c) 2009,2014 Google Inc. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/google/uuid@v1.6.0/dce.go (source SHA-256: `6a4b89e7f6165a673e76f00bebaabf8d0d7eb4356a672c6cb771d95a45fb2695`)
- github.com/google/uuid@v1.6.0/doc.go (source SHA-256: `0b286e2af61ede9f58d2455781ea254284c477b7fce2447aab3203992cc51f9f`)
- github.com/google/uuid@v1.6.0/hash.go (source SHA-256: `afe975c3f3e8b9a972c66165f06e978cc5c2f8ea5a4809dae4c46b4cfebbadce`)
- github.com/google/uuid@v1.6.0/marshal.go (source SHA-256: `933eb226d049268e17809fd37f777c650b9ec8cb3fa988534b792610a64dcc43`)
- github.com/google/uuid@v1.6.0/node.go (source SHA-256: `73a93d192ca24078e0f93b8c1675e13adb442421252d3974d697b0821c4832cf`)
- github.com/google/uuid@v1.6.0/sql.go (source SHA-256: `29281b11e8f97f1ddf73d087854cb401ba42a87c4f59ce45fc69763a4d927f2f`)
- github.com/google/uuid@v1.6.0/time.go (source SHA-256: `e07999a07de5b667dd1dd5792b544ea933e7e839eca03b3469527f86a8bc2881`)
- github.com/google/uuid@v1.6.0/util.go (source SHA-256: `755f50afafb13dd89ae6177da7eb6dc31ae650ece2f9ba52821021e9f8de044e`)
- github.com/google/uuid@v1.6.0/version1.go (source SHA-256: `38a2c229a2aad4bb72d45529304f210b1f31e8ef60bbd2605690a10f0d8a8242`)
- github.com/google/uuid@v1.6.0/version4.go (source SHA-256: `f252aeb4028659d83cbf7b037d4524f7c9b76cde1fdca1f6f7de310dab6f4dcd`)

````text
// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- github.com/google/uuid@v1.6.0/node_net.go (source SHA-256: `78e274b465de8da8690f4ab06281f41025b142c1d093e337fda1d6726b35662b`)

````text
// Copyright 2017 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- github.com/google/uuid@v1.6.0/uuid.go (source SHA-256: `0edec8e34c6b6fe0db31b71a29069a09ed832e3fd04ee0175916b58f2b60e5c1`)

````text
// Copyright 2018 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- github.com/google/uuid@v1.6.0/null.go (source SHA-256: `791638f42d54fcf07ac752df6b2df92593f7da308f944465e2cda443dc3346cf`)

````text
// Copyright 2021 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- github.com/google/uuid@v1.6.0/version6.go (source SHA-256: `013cc9f4a9e606746d38d53f8e651b9222f5042fde567484588ec1d7a0c54905`)
- github.com/google/uuid@v1.6.0/version7.go (source SHA-256: `8002748cc22a032c3782b42f7653a947189b8566c50d246b4a3beb2e075f8d62`)

````text
// Copyright 2023 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## github.com/googleapis/enterprise-certificate-proxy@v0.3.23

Source: https://pkg.go.dev/github.com/googleapis/enterprise-certificate-proxy@v0.3.23

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:J+VGF1wjAM2dU2kKUXZxJdDppRLReWP7KWcclfngQhc=`

### github.com/googleapis/enterprise-certificate-proxy@v0.3.23/LICENSE

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- github.com/googleapis/enterprise-certificate-proxy@v0.3.23/client/client.go (source SHA-256: `884e3301c60a85d847ffdfcd41d20be0c1e9e3a5f0eb66795f25a81b6b1ec749`)
- github.com/googleapis/enterprise-certificate-proxy@v0.3.23/client/util/util.go (source SHA-256: `1032f6c75278844b9b5b1337cf1d2c551e5f274079ba02ca85d2a13090c478a7`)

````text
// Copyright 2022 Google LLC.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## github.com/googleapis/gax-go/v2@v2.26.2

Source: https://pkg.go.dev/github.com/googleapis/gax-go/v2@v2.26.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:ydkmNXxj7bEmmeK5AihkKnWxyOyBR9TDebvp5L5izk8=`

### github.com/googleapis/gax-go/v2@v2.26.2/LICENSE

SHA-256: `b95218cd9607855a6536384c0262922b30a0c2bf56e4ced790240f3a3bac4722`

````text
Copyright 2016, Google Inc.
All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/call_option.go (source SHA-256: `10d58ac1d06eaf2bfe1603217f8626d3c024d13240cc056a66b16b47c27700c7`)
- github.com/googleapis/gax-go/v2@v2.26.2/gax.go (source SHA-256: `0705b0d2fdfb5cf7d69d05bf92e2a235d1bfdfe0009914e790ed9e971c9ea28c`)
- github.com/googleapis/gax-go/v2@v2.26.2/invoke.go (source SHA-256: `37adfb44cd53b98f5f13e8b6cd09b915b79cf20ae4757bfafbb792a26e5d7e03`)

````text
// Copyright 2016, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/header.go (source SHA-256: `c8c1cfea04705cc810d2cfec495860f7edf2604f329d553356a56bee36f245df`)

````text
// Copyright 2018, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/apierror/internal/proto/error.pb.go (source SHA-256: `9818652813df68720ba498a3e4731ff7eeadb1ca94b4df4e5d2b28df9b03faa3`)

````text
// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/apierror/apierror.go (source SHA-256: `4643846951f65c9891c11162f8e9b653a4e1078a9be65edffb719fa7bb32cfea`)

````text
// Copyright 2021, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/apierror/internal/proto/custom_error.pb.go (source SHA-256: `13947a9fb4a30ee8cfa0e81badc9bd4183e3f4f5cd02b905d512eb725f34c06e`)

````text
// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/content_type.go (source SHA-256: `841ea8b61906d504e51f209044efa64bc931c24a7da403a486098268ab011c9d`)
- github.com/googleapis/gax-go/v2@v2.26.2/proto_json_stream.go (source SHA-256: `235096e31fa406b7eaf7d627c46c3a1382fdf3e62f0633320716891ec48b0f7d`)

````text
// Copyright 2022, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/callctx/callctx.go (source SHA-256: `b88fb21ce1b9af2545d48c1c62d6cda15e87fb566819a50b723bb9c8df29b02f`)

````text
// Copyright 2023, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/internallog/internal/internal.go (source SHA-256: `0daab10cc6b0d9d9a595abd33de85204ffbf91ea5c78cfe86e05c8e2fd26dfa6`)
- github.com/googleapis/gax-go/v2@v2.26.2/internallog/internallog.go (source SHA-256: `47064556918c3d84718d479e06a72cc22164c1bf186b0a1bf7595c0df5b4e60d`)

````text
// Copyright 2024, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/feature.go (source SHA-256: `0c9224a2293f069ef4d63f072bf5429ba263aef12aea5b7b28a16ab90cfc7c14`)

````text
// Copyright 2025, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/internal/version.go (source SHA-256: `97d0cce0a1af8e1eba7b282a165c554582bf428113b317360a833c88572802e6`)

````text
// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````

### Original source notice

- github.com/googleapis/gax-go/v2@v2.26.2/telemetry.go (source SHA-256: `fe6ed3bb5de27fb6b50887d2653d9a7090594967977709439a51aa23bad0f8e1`)

````text
// Copyright 2026, Google Inc.
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/hirochachacha/go-smb2@v1.1.0

Source: https://pkg.go.dev/github.com/hirochachacha/go-smb2@v1.1.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:b6hs9qKIql9eVXAiN0M2wSFY5xnhbHAQoCwRKbaRTZI=`

### github.com/hirochachacha/go-smb2@v1.1.0/LICENSE

SHA-256: `8b532d93b35ab59ae9b9ccc08aae66f003c1449b1f7a2d12f5900b008dda1415`

````text
Copyright (c) 2016 Hiroshi Ioka. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/hirochachacha/go-smb2@v1.1.0/internal/crypto/cmac/cmac.go (source SHA-256: `d645c1ac68c24712c9688331d4b18f9831f9a0a637319ecb0c93f897233f7200`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Portions Copyright 2016 Hiroshi Ioka. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/hirochachacha/go-smb2@v1.1.0/internal/crypto/ccm/util.go (source SHA-256: `d0061260ab0f07f82207656b0915c1e24944641e95729ca03e2f31d9ea3829bf`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Portions Copyright 2016 Hiroshi Ioka. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/hirochachacha/go-smb2@v1.1.0/all.go (source SHA-256: `caf55aa970b9812089bd9a1d51591933d97f590e7b965ac9e80ab7486b4c8a36`)

````text
// Original: src/os/path.go
//
// Copyright 2009 The Go Authors. All rights reserved.
// Portions Copyright 2016 Hiroshi Ioka. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/hirochachacha/go-smb2@v1.1.0/filepath.go (source SHA-256: `29a5a7dee6219efa6c48a859194e0bdfa541e52e0d41c175a0af356e3f6e3033`)

````text
// Original: src/path/filepath/match.go
//
// Copyright 2010 The Go Authors. All rights reserved.
// Portions Copyright 2021 Hiroshi Ioka. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/inconshreveable/mousetrap@v1.1.0

Source: https://pkg.go.dev/github.com/inconshreveable/mousetrap@v1.1.0

Targets: windows/amd64

Module checksum: `h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=`

### github.com/inconshreveable/mousetrap@v1.1.0/LICENSE

SHA-256: `4819582701f150b28a563a6cd8ed0bf5a754e3c67b90ad38d78ba4131bf77795`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright 2022 Alan Shreve (@inconshreveable)

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````


## github.com/jeandeaual/go-locale@v0.0.0-20250612000132-0ef82f21eade

Source: https://pkg.go.dev/github.com/jeandeaual/go-locale@v0.0.0-20250612000132-0ef82f21eade

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:FmusiCI1wHw+XQbvL9M+1r/C3SPqKrmBaIOYwVfQoDE=`

### github.com/jeandeaual/go-locale@v0.0.0-20250612000132-0ef82f21eade/LICENSE

SHA-256: `f5a169b17b7b5bc97a3909f190089256dc62adb86eefc65454b0fdca37861c20`

````text
MIT License

Copyright (c) 2020 Alexis Jeandeau

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25

Source: https://pkg.go.dev/github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:YLvr1eE6cdCqjOe972w/cYF+FjW34v27+9Vo5106B4M=`

### github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25/COPYING.txt

SHA-256: `027529c1f33b6ef9280c5b9569e120734dbbee0e5e37436fb9a1bc7a0153b0ae`

````text
Copyright (c) 2012-2015 Jason Summers

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````

### Original source notice

- github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25/reader.go (source SHA-256: `725f93d39ad11e87d9bc284899426ab335374dc67d709d5ddfa36120d782ff36`)

````text
// ◄◄◄ gobmp/reader.go ►►►
// Copyright © 2012 Jason Summers
// Use of this code is governed by an MIT-style license that can
// be found in the readme.md file.
//
// BMP file decoder
//
````

### Original source notice

- github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25/rle.go (source SHA-256: `3e526b563727977628512b3d7a0e60aa7e2147573a7f26af2b4258fedd911800`)

````text
// ◄◄◄ gobmp/rle.go ►►►
// Copyright © 2012 Jason Summers
// Use of this code is governed by an MIT-style license that can
// be found in the readme.md file.
//
// BMP RLE decoder
//
````

### Original source notice

- github.com/jsummers/gobmp@v0.0.0-20230614200233-a9de23ed2e25/writer.go (source SHA-256: `2977b82c4bf54fd1e4cdc186b588220a50ca2a88752d630d2f54b4943759b479`)

````text
// ◄◄◄ gobmp/writer.go ►►►
// Copyright © 2012 Jason Summers
// Use of this code is governed by an MIT-style license that can
// be found in the readme.md file.
//
// BMP file encoder
//
````


## github.com/kr/fs@v0.1.0

Source: https://pkg.go.dev/github.com/kr/fs@v0.1.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:Jskdu9ieNAYnjxsi0LbQp1ulIKZV1LAFgK1tWhpZgl8=`

### github.com/kr/fs@v0.1.0/LICENSE

SHA-256: `dd26a7abddd02e2d0aba97805b31f248ef7835d9e10da289b22e3b8ab78b324d`

````text
Copyright (c) 2012 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/mattn/go-runewidth@v0.0.24

Source: https://pkg.go.dev/github.com/mattn/go-runewidth@v0.0.24

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:cpokDiIn0MGnhdHwuWnJBITySJ20QyNGnY2kR/ay2DU=`

### github.com/mattn/go-runewidth@v0.0.24/LICENSE

SHA-256: `88a2379b3ca34bf5c57127aff9dcb802bbb60ece0805cdbda65b3bd115f971d9`

````text
The MIT License (MIT)

Copyright (c) 2016 Yasuhiro Matsumoto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/mitchellh/colorstring@v0.0.0-20190213212951-d06e56a500db

Source: https://pkg.go.dev/github.com/mitchellh/colorstring@v0.0.0-20190213212951-d06e56a500db

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:62I3jR2EmQ4l5rM/4FEfDWcRD+abF5XlKShorW5LRoQ=`

### github.com/mitchellh/colorstring@v0.0.0-20190213212951-d06e56a500db/LICENSE

SHA-256: `3c377fad2e5ae1d7081c7c2f16da867a87cca1d1f5f1aa7ed0b8a16bb553142a`

````text
The MIT License (MIT)

Copyright (c) 2014 Mitchell Hashimoto

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````


## github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646

Source: https://pkg.go.dev/github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:zYyBkD/k9seD2A7fsi6Oo2LfFZAehjjQMERAvZLEDnQ=`

### github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/LICENSE

SHA-256: `7b850692b15c71706bc6906b51d6375c1404c2d26233078335b6aaa7def8b3f4`

````text
Copyright (c) 2012, Jan Schlicht <jan.schlicht@gmail.com>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
````

### Original source notice

- github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/converter.go (source SHA-256: `6df51a3d43bb310c8913e8aaa203571c2e37f323d3aadbeca75b9256d793100c`)
- github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/filters.go (source SHA-256: `d777e872939cab018547f379b30df65f142ed64148093b67120e17314b2c1b8b`)
- github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/resize.go (source SHA-256: `743e311473ade8b0a0e5505dab7a4082978c71fc961b9537ba78c42118214aa3`)
- github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/thumbnail.go (source SHA-256: `6a21ade4daef5262a5187f9501217358d6f6733e82006df6ceb54b607bfe67b2`)

````text
/*
Copyright (c) 2012, Jan Schlicht <jan.schlicht@gmail.com>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
*/
````

### Original source notice

- github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/nearest.go (source SHA-256: `6d0d16db4a2104f179575cc3d11bb6fb6953f547359b9f0bba83c3c7b1df9d5e`)
- github.com/nfnt/resize@v0.0.0-20180221191011-83c6a9932646/ycc.go (source SHA-256: `2e080b59520cd6bea0a0f0a379c53fcdb1862ca78f818e146a3843a8fcc9a800`)

````text
/*
Copyright (c) 2014, Charlie Vieth <charlie.vieth@gmail.com>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted, provided that the above copyright notice
and this permission notice appear in all copies.

THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES WITH
REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF MERCHANTABILITY AND
FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY SPECIAL, DIRECT,
INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES WHATSOEVER RESULTING FROM LOSS
OF USE, DATA OR PROFITS, WHETHER IN AN ACTION OF CONTRACT, NEGLIGENCE OR OTHER
TORTIOUS ACTION, ARISING OUT OF OR IN CONNECTION WITH THE USE OR PERFORMANCE OF
THIS SOFTWARE.
*/
````


## github.com/nicksnyder/go-i18n/v2@v2.5.1

Source: https://pkg.go.dev/github.com/nicksnyder/go-i18n/v2@v2.5.1

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:IxtPxYsR9Gp60cGXjfuR/llTqV8aYMsC472zD0D1vHk=`

### github.com/nicksnyder/go-i18n/v2@v2.5.1/LICENSE

SHA-256: `8fb3aead091abd1579670cd4d805c7923b66ee97b771148f4ef928a65c259328`

````text
Copyright (c) 2014 Nick Snyder https://github.com/nicksnyder

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
````


## github.com/pkg/browser@v0.0.0-20240102092130-5ac0b6a4141c

Source: https://pkg.go.dev/github.com/pkg/browser@v0.0.0-20240102092130-5ac0b6a4141c

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:+mdjkGKdHQG3305AYmdv1U2eRNDiU2ErMBj1gwrq8eQ=`

### github.com/pkg/browser@v0.0.0-20240102092130-5ac0b6a4141c/LICENSE

SHA-256: `4801bb0c5d70b980b3e1f63d1d2df8d84fec22548e3d73f03a68df1fb808791f`

````text
Copyright (c) 2014, Dave Cheney <dave@cheney.net>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/pkg/sftp@v1.13.12

Source: https://pkg.go.dev/github.com/pkg/sftp@v1.13.12

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:Myhzy/ZzLPnCtE2F2WtDHsgs8RfgMnH2gfbHuCJhkMU=`

### github.com/pkg/sftp@v1.13.12/LICENSE

SHA-256: `f5b5087b8e596d015cfea7db26527fe67a196c18ca37415c288dd2534ff74c37`

````text
Copyright (c) 2013, Dave Cheney
All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

 * Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 * Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````


## github.com/rivo/uniseg@v0.4.7

Source: https://pkg.go.dev/github.com/rivo/uniseg@v0.4.7

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:WUdvkW8uEhrYfLC4ZzdpI2ztxP1I582+49Oc5Mq64VQ=`

### github.com/rivo/uniseg@v0.4.7/LICENSE.txt

SHA-256: `a59885f5f0f3b3c07cf9444db5fe399b6f0791eac82055f7b85cc65500551039`

````text
MIT License

Copyright (c) 2019 Oliver Kuederle

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### Original source notice

- github.com/rivo/uniseg@v0.4.7/lineproperties.go (source SHA-256: `d73432ab989195312fdaa4f4dcea34c7bfb478160b98ad26fbc31df94f898b45`)

````text
//         COPYRIGHT SIGN
````

### Original source notice

- github.com/rivo/uniseg@v0.4.7/lineproperties.go (source SHA-256: `d73432ab989195312fdaa4f4dcea34c7bfb478160b98ad26fbc31df94f898b45`)

````text
//         SOUND RECORDING COPYRIGHT
````

### Original source notice

- github.com/rivo/uniseg@v0.4.7/graphemeproperties.go (source SHA-256: `321483625d9a4aa92e4a5730df51c59d5b49725cd20e3879749488e5628b3da8`)
- github.com/rivo/uniseg@v0.4.7/wordproperties.go (source SHA-256: `10136e6c3c8ea619276428d59278e312aa4ce33b2e50d98f3a618f162b3a5864`)

````text
// E0.6   [1] (©️)       copyright
````

### Original source notice

- github.com/rivo/uniseg@v0.4.7/eastasianwidth.go (source SHA-256: `1cd2368ddedc56182d2f64576f2a713a2ad6234265d940575285f61006527c5b`)

````text
// So         COPYRIGHT SIGN
````

### Original source notice

- github.com/rivo/uniseg@v0.4.7/eastasianwidth.go (source SHA-256: `1cd2368ddedc56182d2f64576f2a713a2ad6234265d940575285f61006527c5b`)

````text
// So         SOUND RECORDING COPYRIGHT
````


## github.com/rymdport/portal@v0.4.2

Source: https://pkg.go.dev/github.com/rymdport/portal@v0.4.2

Targets: linux/amd64

Module checksum: `h1:7jKRSemwlTyVHHrTGgQg7gmNPJs88xkbKcIL3NlcmSU=`

### github.com/rymdport/portal@v0.4.2/LICENSE

SHA-256: `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````


## github.com/schollz/progressbar/v3@v3.19.1

Source: https://pkg.go.dev/github.com/schollz/progressbar/v3@v3.19.1

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:iv8BgwOvdML/S3p84uBpy/IMigv4U9594vPZYa2EdrU=`

### github.com/schollz/progressbar/v3@v3.19.1/LICENSE

SHA-256: `0c27314bc4294775c1c7768690fa90b1884a3e27b72b92376ee505963f916e8a`

````text
MIT License

Copyright (c) 2017 Zack

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/spf13/cobra@v1.10.2

Source: https://pkg.go.dev/github.com/spf13/cobra@v1.10.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:DMTTonx5m65Ic0GOoRY2c16WCbHxOOw6xxezuLaBpcU=`

### github.com/spf13/cobra@v1.10.2/LICENSE.txt

SHA-256: `5e3400b93bbb099e83e52bab885e7441750673c21f97988ca3f1240639b63283`

````text
                                Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.
````

### Original source notice

- github.com/spf13/cobra@v1.10.2/active_help.go (source SHA-256: `92bde33086d33e413b9ac945ddb38f7487f9ce91067453a7ce97fd62fa59cdc5`)
- github.com/spf13/cobra@v1.10.2/args.go (source SHA-256: `db24bf6cce3df231100f6b79387e8b4f96b994b438a05d653b0c0a869f605660`)
- github.com/spf13/cobra@v1.10.2/bash_completions.go (source SHA-256: `7370190d5344eec0325b2fd2dd7a0f39ca2acc5881f818adb5124b16ee7ac507`)
- github.com/spf13/cobra@v1.10.2/bash_completionsV2.go (source SHA-256: `1a90ae94435505c55c315c5e40a64ee3fbcf4aa790f741ecbeff5ee37116f149`)
- github.com/spf13/cobra@v1.10.2/cobra.go (source SHA-256: `aae7502e789e42469bcb2fc110af87fd02451d3b509c5b0bed782ada17338c5f`)
- github.com/spf13/cobra@v1.10.2/command.go (source SHA-256: `f79613721ad04c9b438dfb30fe0ac3deca85c487658b0aa04a1ae6e4395bfd90`)
- github.com/spf13/cobra@v1.10.2/command_notwin.go (source SHA-256: `505a8334864917e8439b3012e8e47d880643793bf5741bc048884a85037a8ac2`)
- github.com/spf13/cobra@v1.10.2/command_win.go (source SHA-256: `a9f85ba76b54df4b2589546562b325c7b1225388d5b177214d4b193516556390`)
- github.com/spf13/cobra@v1.10.2/completions.go (source SHA-256: `7d94de34c9ebf37da817a1dec576e672c41c4bf7d1e42140f39e7ad36bee9eee`)
- github.com/spf13/cobra@v1.10.2/fish_completions.go (source SHA-256: `86ee8a7c124ab4034290ec156f32ac0c051038fd560bb8d3c95214bdaf2e858a`)
- github.com/spf13/cobra@v1.10.2/flag_groups.go (source SHA-256: `3b0572e1fc2fc6e8ef62319b4294cdda86259fb0b1afebc18b03fa643e431251`)
- github.com/spf13/cobra@v1.10.2/powershell_completions.go (source SHA-256: `6357a6d1adc9e798f21731cded574fc3e6e95240c92ee471541d61478b63351a`)
- github.com/spf13/cobra@v1.10.2/shell_completions.go (source SHA-256: `38244d5514d8de4c7bd1e03890992a9ee73885af53aab4294d8bed2d76f8938d`)
- github.com/spf13/cobra@v1.10.2/zsh_completions.go (source SHA-256: `a80998389020b259308582d539cb3ae0979f83b55e8ffa68b20c46b015feae24`)

````text
// Copyright 2013-2023 The Cobra Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## github.com/spf13/pflag@v1.0.9

Source: https://pkg.go.dev/github.com/spf13/pflag@v1.0.9

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:9exaQaMOCwffKiiiYk6/BndUBv+iRViNW+4lEMi0PvY=`

### github.com/spf13/pflag@v1.0.9/LICENSE

SHA-256: `b8514c577c1c4b46cee454d5a882b15fa411e72c5bd7f801f241591789fce61a`

````text
Copyright (c) 2012 Alex Ogier. All rights reserved.
Copyright (c) 2012 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/spf13/pflag@v1.0.9/flag.go (source SHA-256: `5699e18dd47610cf862dca1bb046f5d00a652fd412d9ac2c1a0f79b9c691c8cb`)
- github.com/spf13/pflag@v1.0.9/golangflag.go (source SHA-256: `adba9f0edd9f0d6db4517aa63f666a32a7840a9d95506f9998d53e37d3f67cf9`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c

Source: https://pkg.go.dev/github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:km8GpoQut05eY3GiYWEedbTT0qnSxrCjsVbb7yKY1KE=`

### github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/LICENSE

SHA-256: `cda28abc46a25d34ad876ec119ad3c29edbb8634d5fef46ccf2e1d63acd9ca34`

````text
BSD 3-Clause License

Copyright (c) 2018, Steven R Wiley
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/testdata/LICENSE

SHA-256: `c654732ce4d4761102c29a97006abf63f2c2294c1558dd73d6acca5a327840f8`

````text
Test icons in the testdata folder are from either www.freepik.com or www.flaticon.com. They request any wedsite using these icons place the following in thei web pages:

<div>Icons made by <a href="http://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>

<div>Icons made by <a href="https://www.flaticon.com/authors/flat-icons" title="Flat Icons">Flat Icons</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a> is licensed by <a href="http://creativecommons.org/licenses/by/3.0/" title="Creative Commons BY 3.0" target="_blank">CC 3.0 BY</a></div>
````

### Original source notice

- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/definitions.go (source SHA-256: `49b2ad5bb7098ab5a6ea5618e7b221b25132a78807778cbdad133bf5f6e1be88`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/draw.go (source SHA-256: `08127f9847491041f748b4b69d64b60451731599a70bb3a8c61d634d1334ddc8`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/icon_cursor.go (source SHA-256: `df83e208e1f86e337813d247eac026d7eb49e84a0b30c7142dad360f8bee36bd`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/path_cursor.go (source SHA-256: `24e2ed3ca4b375854d130cbe96535b0d5f7c254e90cc14cbb5b517ebe4414236`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/path_style.go (source SHA-256: `e61be3deac20e1fc1f7d57cfd6a40ecf7fe8ff53ac2eaef9bb49e7081802ed30`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/public.go (source SHA-256: `479f4d6bc93d144f775e85ab4fa73e081bb1ca0840471c3c0611ffcd3dbfe987`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/svg_icon.go (source SHA-256: `087831bf805003e9793c101680988cc6447303423b957880f0d36053a8f207e5`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/svg_path.go (source SHA-256: `934ec3bfb250e4377eacba11e8325dc83e7396fb7f0992a8a1a9e2a2b75de856`)
- github.com/srwiley/oksvg@v0.0.0-20221011165216-be6e8873101c/utils.go (source SHA-256: `ace28872b05986682e94dc030eb21406e35d07ca10ae1acd8f5bf3753e5df54c`)

````text
// Copyright 2017 The oksvg Authors. All rights reserved.
// created: 2/12/2017 by S.R.Wiley
//
// utils.go implements translation of an SVG2.0 path into a rasterx Path.
````


## github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef

Source: https://pkg.go.dev/github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:Ch6Q+AZUxDBCVqdkI8FSpFyZDtCVBc2VmejdNrm5rRQ=`

### github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/LICENSE

SHA-256: `cda28abc46a25d34ad876ec119ad3c29edbb8634d5fef46ccf2e1d63acd9ca34`

````text
BSD 3-Clause License

Copyright (c) 2018, Steven R Wiley
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/stroke.go (source SHA-256: `f595d7cc9a6148c92c22efb5bc0b71476c01b5dfc804ccd9b0af71782f1f97c2`)

````text
// Copyright 2017 by the rasterx Authors. All rights reserved.
//
// created: 2017 by S.R.Wiley
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/dash.go (source SHA-256: `20061e9f08ea445ad5bfc97fa2c4d3eb09587995f3d3ec67af64836be10c5e04`)

````text
// Copyright 2017 by the rasterx Authors. All rights reserved.
//_
// created: 2017 by S.R.Wiley
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/fill.go (source SHA-256: `102c1b4c27671309a6bc67a676e0b5237577a2c036104ced4be5fbbb4fe9ab1e`)

````text
// Copyright 2018 by the rasterx Authors. All rights reserved.
//_
// Created 2017 by S.R.Wiley
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/shapes.go (source SHA-256: `bf7f9b2fe77e6fa0f72071229a5f82f615cfc998234c5ac5278cb013837e52f5`)

````text
// Copyright 2018 by the rasterx Authors. All rights reserved.
//_
// created: 2/06/2018 by S.R.Wiley
// Functions that rasterize common shapes easily.
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/gradient.go (source SHA-256: `cb2c68cbb2d825765c8347c015da1c17b48ed0d5bc8796c70515aedb95ba9ca7`)

````text
// Gradient implementation fo rasterx package
// Copyright 2018 All rights reserved.
// Created: 5/12/2018 by S.R.Wiley
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/matrix.go (source SHA-256: `c21bec1920f799527e6f5598c4be447e7a5bd78fbb50128ff4dce6ff4d51f13d`)

````text
// Implements SVG style matrix transformations.
// https://developer.mozilla.org/en-US/docs/Web/SVG/Attribute/transform
// Copyright 2018 All rights reserved.
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/scan.go (source SHA-256: `cb80cc0c1267119fd3f97c1d85066a06aecaf9b005a3c4cf5b13e560c73b0c05`)

````text
// Package rasterx implements a rasterizer in go.
// By default rasterx uses ScannerGV to render images
// which uses the rasterizer in the golang.org/x/image/vector package.
// The freetype rasterizer under the GNU license can also be used, by
// downloading the scanFT package.
//
// Copyright 2018 All rights reserved.
// Created: 5/12/2018 by S.R.Wiley
````

### Original source notice

- github.com/srwiley/rasterx@v0.0.0-20220730225603-2ab79fcdd4ef/geomx.go (source SHA-256: `0ae1c97df71b754ef40034e273192332566aac979f1afd638284665cd7d6fbe2`)

````text
// geomx adds some geometry functions needed by rasterx
// Copyright 2017 by the rasterx Authors. All rights reserved.
// Created: 2/12/2017 by S.R.Wiley
````


## github.com/yuin/goldmark@v1.8.2

Source: https://pkg.go.dev/github.com/yuin/goldmark@v1.8.2

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:kEGpgqJXdgbkhcOgBxkC0X0PmoPG1ZyoZ117rDVp4zE=`

### github.com/yuin/goldmark@v1.8.2/LICENSE

SHA-256: `c1d6f653c38c80294dc7994647ba0dcc508ce64a86996303c804efa2f7982bff`

````text
MIT License

Copyright (c) 2019 Yusuke Inuzuka

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````


## github.com/zalando/go-keyring@v0.2.8

Source: https://pkg.go.dev/github.com/zalando/go-keyring@v0.2.8

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:6sD/Ucpl7jNq10rM2pgqTs0sZ9V3qMrqfIIy5YPccHs=`

### github.com/zalando/go-keyring@v0.2.8/LICENSE

SHA-256: `14bc00f1a298eb6af250216b9e60027a4386c8a84fa2c7ddc2b5c10908e06d0d`

````text
The MIT License (MIT)

Copyright (c) 2016 Zalando SE

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### github.com/zalando/go-keyring@v0.2.8/internal/shellescape/LICENSE

SHA-256: `01a4609ac3c7ac0a26d65a967c26c0deb146aa01bdcc8531048a70be52695556`

````text
The MIT License (MIT)

Copyright (c) 2016 Alessio Treglia

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
````

### Original source notice

- github.com/zalando/go-keyring@v0.2.8/keyring_darwin.go (source SHA-256: `8fd983f3af38bfd64b62949788b0ad5f7e03d24c87fde9a994b1e1a582ba2792`)

````text
// Copyright 2013 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## go.opentelemetry.io/auto/sdk@v1.2.1

Source: https://pkg.go.dev/go.opentelemetry.io/auto/sdk@v1.2.1

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:jXsnJ4Lmnqd11kwkBV2LgLoFMZKizbCi5fNZ/ipaZ64=`

### go.opentelemetry.io/auto/sdk@v1.2.1/LICENSE

SHA-256: `c71d239df91726fc519c6eb72d318ec65820627232b2f796219e87dcf35d0ab4`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- go.opentelemetry.io/auto/sdk@v1.2.1/doc.go (source SHA-256: `26a636492b7dae83b65b5fe223a03fce8abdc7f4fe49e3646e6b5b8968abc18a`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/attr.go (source SHA-256: `8ee1b06c52ec9edcd60f6722f2ccb1da60401619a2c26081d06ea55a757fb558`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/doc.go (source SHA-256: `faf0de2a44e30b72bb1577e4a77fa0957dfaa123c5dfafe601ba3c43e85fb9d8`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/id.go (source SHA-256: `d5ea09864e3ef0ffeb0c5f4930bc9266834ce68085d972ee096f6991e9055064`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/number.go (source SHA-256: `d2c4356084e9f1dcb62e5dcc757f9eb5088b6f4a7b72434db37b8359a8add36e`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/resource.go (source SHA-256: `5eae9f7709acc4255ed376f2fa02ebc60925ce20ed178d968db80d9d846e9adc`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/scope.go (source SHA-256: `4dea63d4152ace4fdeda335dc1732c514779f2d28ac54c59ca575965c75cf2e8`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/span.go (source SHA-256: `065834dcff0df37cc689d87f4fbbbc9589c0160374e3ce3a85407f070de2fa6d`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/status.go (source SHA-256: `db31f8b262c2adbfd19b0bbc13b52c8b8ad98d82460e852a2154b372ccaa5d41`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/traces.go (source SHA-256: `fc4e9f31180a606a8558a6cfa21b9197bdae2f6c510eae9a7116bd69c62bc15f`)
- go.opentelemetry.io/auto/sdk@v1.2.1/internal/telemetry/value.go (source SHA-256: `7cfb5acb61b2a7cc581f93c07cd7ec394c8c456490e475129a4a4ade4c14fc9a`)
- go.opentelemetry.io/auto/sdk@v1.2.1/limit.go (source SHA-256: `e3befa89f1059bcef48189f437f74e9f5f1f3dbaa3c44421bbdde368594f03e0`)
- go.opentelemetry.io/auto/sdk@v1.2.1/span.go (source SHA-256: `db589f08a63d1ffc5b1d3505512221977c6719433d9c58c37545e75ffc375bc7`)
- go.opentelemetry.io/auto/sdk@v1.2.1/tracer.go (source SHA-256: `6418ba2ae2cf0d16725a81a5ed5f404faa630a72da457afb8cb16d9ee6f431c9`)
- go.opentelemetry.io/auto/sdk@v1.2.1/tracer_provider.go (source SHA-256: `0abfec6c21fe9a9e1d560c2ec982852fa9340fd719b1787111eb7347ca2ea4dc`)

````text
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
````


## go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0

Source: https://pkg.go.dev/go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:8tvICD4vSTOOsNrsI4Ljf6C+6UKvpTEH5XY3JMoyPoo=`

### go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/LICENSE

SHA-256: `1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

--------------------------------------------------------------------------------

Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/common.go (source SHA-256: `7e2c747d5fc33067062aa60be13d26e69fba7a49720808be533930b66390c0bb`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/config.go (source SHA-256: `a04ba07bfcfda7dfa6e90e3c560a7cf6a4085c38f36d427d8df8ebd3605a5fa3`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/doc.go (source SHA-256: `fb774dc4b9f20867b2eb1dc43afe836b5d42313b702b09c5d77968a775cd05c5`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/handler.go (source SHA-256: `1dd0492525e6cb4c83d0cebb2239354368347628e0d12340ecfa0dc069e18aaf`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/request/body_wrapper.go (source SHA-256: `29b864f525f69b2cb91c1362cf51c580cc042a4c9dbc02ddf54978d51678496c`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/request/gen.go (source SHA-256: `eee3327df0c013146d72e7d2d3f2ec45a3a21599b92a14e17d97106fb7a05121`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/request/resp_writer_wrapper.go (source SHA-256: `38958f42437388f51c7c883c4f69ae338e35b7fe4247577256317e8af6f30de3`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/semconv/client.go (source SHA-256: `89d01ac6c8b86d133a2723932fd2286cad96fe1d852f1aace17270e2c117f9bd`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/semconv/gen.go (source SHA-256: `b701e83f3861744ef3bda3e924d307a2b8f46dc6d0cd3980317644d374bfdd97`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/semconv/server.go (source SHA-256: `2a9abdd42f43a360fcb60fabd41410bfb2f50efbd2b888e75f4f0b4118247ce1`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/internal/semconv/util.go (source SHA-256: `ae86beb8295e6fe09c56edb19f500f2db7c71d4779195044be86f8f7c5803995`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/labeler.go (source SHA-256: `d4d98d5ad6e9b496c9c6f5e16a75ece96074fe163d491d04c90156a603303c2b`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/start_time_context.go (source SHA-256: `311e76549fb72f891ba096bb3fd69c385a1bc6ba187fc4b4903f9ad10da6879d`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/transport.go (source SHA-256: `8da58a37595f208c3227fe7037875b91d7552f9942580b2e7986544734dab541`)
- go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp@v0.69.0/version.go (source SHA-256: `15c2954cb0b7e2ef2a0fbeb8379a46615c1ebad2032eaf461769f2ac9127e1b8`)

````text
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
````


## go.opentelemetry.io/otel/metric@v1.45.0

Source: https://pkg.go.dev/go.opentelemetry.io/otel/metric@v1.45.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:7Eg1uH7CJ5cXv9is6tnBe1FI6rj1nwUdbFypRm3br/M=`

### go.opentelemetry.io/otel/metric@v1.45.0/LICENSE

SHA-256: `1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

--------------------------------------------------------------------------------

Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- go.opentelemetry.io/otel/metric@v1.45.0/asyncfloat64.go (source SHA-256: `552e47997273fd75fcd14b1e6e948f6d94153072a4b621e911535979310a9c4c`)
- go.opentelemetry.io/otel/metric@v1.45.0/asyncint64.go (source SHA-256: `07f51ef78fadfaa5f153c429b18934b4ffaeb80fec33145de5e805c88c4e75af`)
- go.opentelemetry.io/otel/metric@v1.45.0/config.go (source SHA-256: `3bb4d8ad5b078afcc7536ecd5e3a38dc466693c122f4c5012f8b0c949a3fc22b`)
- go.opentelemetry.io/otel/metric@v1.45.0/doc.go (source SHA-256: `5f057f47f648b8641ed912d24b5342098a2dad58fc94c4952e6e8ed6a83f9263`)
- go.opentelemetry.io/otel/metric@v1.45.0/embedded/embedded.go (source SHA-256: `8d1170cacbb4392f3cc4ae5606fbeecd76ef1aac2465a0956b92388709e5ab14`)
- go.opentelemetry.io/otel/metric@v1.45.0/instrument.go (source SHA-256: `1fe7c2ccd5e90371e847d501e9bc2c34349634bac9253b008d55d1a2aeca5d84`)
- go.opentelemetry.io/otel/metric@v1.45.0/meter.go (source SHA-256: `dc393dba980c0207d7cd14e47bb15deff4d7a0f413498c3163e27290486a7d74`)
- go.opentelemetry.io/otel/metric@v1.45.0/noop/noop.go (source SHA-256: `08d6b6cb06e3b8601c590f6bcc927c7f7e2fba1f544b50ed25893c023c06cd7c`)
- go.opentelemetry.io/otel/metric@v1.45.0/syncfloat64.go (source SHA-256: `9d4324d520010983919dcb4c0bdd0ff61cdaa887906c984c460e80e0505ed27f`)
- go.opentelemetry.io/otel/metric@v1.45.0/syncint64.go (source SHA-256: `e91c41be06434ba135a1a5ed193f71fee911ac742a0fa4eac4057fe91cbb937f`)

````text
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
````


## go.opentelemetry.io/otel/trace@v1.45.0

Source: https://pkg.go.dev/go.opentelemetry.io/otel/trace@v1.45.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:l/mP6Uv7oNO7/TblbhpbgMidxhq1uO/rPsikOyVhxag=`

### go.opentelemetry.io/otel/trace@v1.45.0/LICENSE

SHA-256: `1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

--------------------------------------------------------------------------------

Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- go.opentelemetry.io/otel/trace@v1.45.0/auto.go (source SHA-256: `68ab9d297278a52b050d3a59ea2ea2c5df90bf04c4c5749724f921737e33e724`)
- go.opentelemetry.io/otel/trace@v1.45.0/config.go (source SHA-256: `c622c7a4ef82972302467f0a9888404431bb63884db13d0d4ce70dd0c3525167`)
- go.opentelemetry.io/otel/trace@v1.45.0/context.go (source SHA-256: `1d915c483364b5ecbaccba8cb96e79421f125482a892ce8a88d0141f484c917b`)
- go.opentelemetry.io/otel/trace@v1.45.0/doc.go (source SHA-256: `1e8208f56ff8741291081bffa1ed81bb252fac4b4602440b70873a0b65dadc19`)
- go.opentelemetry.io/otel/trace@v1.45.0/embedded/embedded.go (source SHA-256: `ac35554d63e2861a68393304f463fc01cc7cf4284e1612921cdf6017a28ef77d`)
- go.opentelemetry.io/otel/trace@v1.45.0/hex.go (source SHA-256: `213aa3afee8ccf8e7c3ed90255170545d8b1af5413cf529aa6b15ffbdb9b2dd5`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/attr.go (source SHA-256: `e7bbf4ac8fea57d0f0940abdcf6aa6d9ceb1d8af32e0f36afa24141dbb73ae06`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/doc.go (source SHA-256: `faf0de2a44e30b72bb1577e4a77fa0957dfaa123c5dfafe601ba3c43e85fb9d8`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/id.go (source SHA-256: `60947ddeb9648cc7e7b6e55be6b840f6eabf0a83a4dad9453360dca47f429911`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/number.go (source SHA-256: `33283d3031a7ed6788c35e4136f29488459e95301f388a0f3e278033fba39f1f`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/resource.go (source SHA-256: `5eae9f7709acc4255ed376f2fa02ebc60925ce20ed178d968db80d9d846e9adc`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/scope.go (source SHA-256: `4dea63d4152ace4fdeda335dc1732c514779f2d28ac54c59ca575965c75cf2e8`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/span.go (source SHA-256: `9c2aa63ac9ec6ddce022cb416eb61a35b8504d3d1c600351d74ea36c7d97223d`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/status.go (source SHA-256: `0da2eefc5e5d6460f030b66ef83f2e16a84f25608f7a260c1af68a476a30ba29`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/traces.go (source SHA-256: `fc4e9f31180a606a8558a6cfa21b9197bdae2f6c510eae9a7116bd69c62bc15f`)
- go.opentelemetry.io/otel/trace@v1.45.0/internal/telemetry/value.go (source SHA-256: `b6f20961e4133d876fddbeaab5b51477f93d5271b941ac540509cb41991d1acf`)
- go.opentelemetry.io/otel/trace@v1.45.0/nonrecording.go (source SHA-256: `1effdb45c4c6aaf6d517d437c1eba06a260df4f7a1fa05ed19cd2793b1724c75`)
- go.opentelemetry.io/otel/trace@v1.45.0/noop.go (source SHA-256: `dbe691da03743b7b2b4dc8ee254ae9e4ce6e3e652407a611c3b93d38a1ac0600`)
- go.opentelemetry.io/otel/trace@v1.45.0/noop/noop.go (source SHA-256: `bd10d239a58a85996933a68213f01bac927778bf525a9881f752d6dac2e0846c`)
- go.opentelemetry.io/otel/trace@v1.45.0/provider.go (source SHA-256: `785e38b85cdec697b94f58c4884d3c93aa2ce81debe57abde0eedb705125dfe6`)
- go.opentelemetry.io/otel/trace@v1.45.0/span.go (source SHA-256: `86772ec971f4bd0e995aee1d98e8f226b5174b5b348795a88f5c2ccbf584b10e`)
- go.opentelemetry.io/otel/trace@v1.45.0/trace.go (source SHA-256: `e40a71ee8aa42d9efaa39f06fb86da8f042a0a5e43cbe46600cfb9ec709b4bcc`)
- go.opentelemetry.io/otel/trace@v1.45.0/tracer.go (source SHA-256: `42ad650fbcc6d4aad56fac781d654300af2f6e09b24b8a848d16ef1062587599`)
- go.opentelemetry.io/otel/trace@v1.45.0/tracestate.go (source SHA-256: `8d0ceb2c09d155ef2c9bf605aac948495a31aab0437055274414dc7e0030d581`)

````text
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
````


## go.opentelemetry.io/otel@v1.45.0

Source: https://pkg.go.dev/go.opentelemetry.io/otel@v1.45.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:pdrWmLHofpubmArBv1LgFSv1Z0Ie/ppdZzu+kUN5EeU=`

### go.opentelemetry.io/otel@v1.45.0/LICENSE

SHA-256: `1ae07514be1d7bb33f0698f8d91fb51b8b9fe1463157ec1c72081a49b9bc6f40`

````text
                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

--------------------------------------------------------------------------------

Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- go.opentelemetry.io/otel@v1.45.0/attribute/doc.go (source SHA-256: `e1efce9865b3a2567e5356f5d055131e51dd6cf7b658bd4effd8b1fc64539165`)
- go.opentelemetry.io/otel@v1.45.0/attribute/encoder.go (source SHA-256: `e2cd9864a3edc3cddcac6702648f934f8a9a3d6ab3f37f0079af034e9be1c295`)
- go.opentelemetry.io/otel@v1.45.0/attribute/filter.go (source SHA-256: `f7acb42faaca96050fa101e49b02fb5c47fa225bf3b55c9c7b48593138aa1721`)
- go.opentelemetry.io/otel@v1.45.0/attribute/hash.go (source SHA-256: `c448244c33ee3ceb237350ff4a984fc8e2bf3b3a26a0f852915673ca642b9fd5`)
- go.opentelemetry.io/otel@v1.45.0/attribute/internal/attribute.go (source SHA-256: `3bbe420e82d8468f05c3b09bb1cdf61a02ffd23592a3e3a9bbfe8318bccc460a`)
- go.opentelemetry.io/otel@v1.45.0/attribute/internal/xxhash/xxhash.go (source SHA-256: `25bced8e5a3222eecb7abd60da41c25a9b245247a1d5ba8aa8d3b6b53e24dfdb`)
- go.opentelemetry.io/otel@v1.45.0/attribute/iterator.go (source SHA-256: `1685c54692348f86a4584a7c5c1656d43a5acf5ad133510c1e98562771ebcfae`)
- go.opentelemetry.io/otel@v1.45.0/attribute/key.go (source SHA-256: `085c984f3c077f0202fe451680ac983f9d0106a4bac372865e68ff7fca52ac28`)
- go.opentelemetry.io/otel@v1.45.0/attribute/kv.go (source SHA-256: `6664127e3449a013c2186a7b85b0b296c851ab3f6c9a35bf0c21f606cc5d5a5b`)
- go.opentelemetry.io/otel@v1.45.0/attribute/rawhelpers.go (source SHA-256: `e1072ec2f132e4aa2e0ea7b73dded46ec2b61d6fff84331081a6a025e5bea2eb`)
- go.opentelemetry.io/otel@v1.45.0/attribute/set.go (source SHA-256: `c9b2a514f862196841690ba4c6aa334e4af8b34cfa7747aeedd56b0618935a7d`)
- go.opentelemetry.io/otel@v1.45.0/attribute/value.go (source SHA-256: `a008a9e6bcc53c0fa599bf80182a7e0e74b262190b76af6a42d55c0bfa89a62d`)
- go.opentelemetry.io/otel@v1.45.0/baggage/baggage.go (source SHA-256: `43f3024d5ccce46ddd5fda961749a4aa8437dd2d73e7f6557d385ed9b99cef13`)
- go.opentelemetry.io/otel@v1.45.0/baggage/context.go (source SHA-256: `53ab84584b05df202d71bac9741cda5b4e8b2946fb358e1852733d50cc95f9a0`)
- go.opentelemetry.io/otel@v1.45.0/baggage/doc.go (source SHA-256: `0a8d489dafd81ed1138bd9c4b0f49e4eede0fa6cfdabedafaf2f8f529b69bf38`)
- go.opentelemetry.io/otel@v1.45.0/codes/codes.go (source SHA-256: `ee3066f71e9295f9ab2e4bedd55061d09de2ca2a0d5822f85dac8faabfffa52e`)
- go.opentelemetry.io/otel@v1.45.0/codes/doc.go (source SHA-256: `80b269087fd663e51cb4630b3d05eb04e7a3f2c17f12c3f34db25d6625129949`)
- go.opentelemetry.io/otel@v1.45.0/doc.go (source SHA-256: `0f3686bfd69e3e96e77190915fa16b3ed3dc226cf2d48f5038749e12e78ccc9c`)
- go.opentelemetry.io/otel@v1.45.0/error_handler.go (source SHA-256: `af0d2c365189797ddbf4ab4c7fbda5712e631d049c780add95f8d5ba02376254`)
- go.opentelemetry.io/otel@v1.45.0/handler.go (source SHA-256: `7174c1748411c8041e83a2baab2251687c2eeac5c4801ee3787b273092376616`)
- go.opentelemetry.io/otel@v1.45.0/internal/baggage/baggage.go (source SHA-256: `033a160c1a6b4bf5bc5c43c845c0ecba0632605d2d0d642260feaf8b19127249`)
- go.opentelemetry.io/otel@v1.45.0/internal/baggage/context.go (source SHA-256: `bf848bef4d075de9bda6290f8755a1dcfa4ed82226e599a8687d5c751e8f2398`)
- go.opentelemetry.io/otel@v1.45.0/internal/errorhandler/errorhandler.go (source SHA-256: `b8cde0b4dfc6cabfdc967cad653596ed3f8245c381c0351e662ea5ed2da7413f`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/handler.go (source SHA-256: `6640e56502009d8bc51b0315ceb11d4f0200f0b7b89516300255ab374a04632c`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/instruments.go (source SHA-256: `3c6a6dd8c775bdb0bc9dbac9d12e506226791c72bbccfc86ad704bacd4570354`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/internal_logging.go (source SHA-256: `5aec503d20272ef97b1baec6328c7f3335aaee12291d7e93dbea5f04d048d150`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/meter.go (source SHA-256: `6793aa02c8ff28f0a0972a67aba1fc92bdb6882ff7a31fe665bb8639152f735c`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/propagator.go (source SHA-256: `d343fbdddd1c69743f2493ea28c65ee029915cfaabb983005e936ec734a41610`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/state.go (source SHA-256: `11ddf161f203558e23476bfd564422da4fd13bb737b6fd9bfe51c97d63e20ae0`)
- go.opentelemetry.io/otel@v1.45.0/internal/global/trace.go (source SHA-256: `207b13705cb37e2bf5f67b1ac2ff8db70ea728cc2cf8ec2fa45bb808a15e9bb3`)
- go.opentelemetry.io/otel@v1.45.0/internal_logging.go (source SHA-256: `dca891ad6a89341ffa3d26ae174d967faa1da6c3ddc4a9d313fd62b7b780baee`)
- go.opentelemetry.io/otel@v1.45.0/metric.go (source SHA-256: `f1a2f91357ab90ebc521d2554265b95202c113d2a81983edd0c0ccba1a80d0e3`)
- go.opentelemetry.io/otel@v1.45.0/propagation.go (source SHA-256: `d0fc816c0247c5bc681483a67052dcb3092638d4d4781866dfe1800a40272f95`)
- go.opentelemetry.io/otel@v1.45.0/propagation/baggage.go (source SHA-256: `fd76dbe1277544a45cb105dd6e06441943183c5ec6898ba8e568b930c9c8f0ae`)
- go.opentelemetry.io/otel@v1.45.0/propagation/doc.go (source SHA-256: `e41de2d6cd4bfa9e3bad549634869cbcf60a05ab7a644cbb768e9c64897e0714`)
- go.opentelemetry.io/otel@v1.45.0/propagation/propagation.go (source SHA-256: `a445a83646b14b366a5b9ed29eb5b4efe9813873473c31b10642ded3d212d05e`)
- go.opentelemetry.io/otel@v1.45.0/propagation/trace_context.go (source SHA-256: `d1f273a92504e95cd45d91fca53eb6eba248209d107bfd2da8c8998790611025`)
- go.opentelemetry.io/otel@v1.45.0/semconv/internal/metricpool/pool.go (source SHA-256: `6efca107fc9ef6316c849fe760df07c0db16f4b5d4af707de3a2e9c930813795`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.37.0/attribute_group.go (source SHA-256: `38a6475feb50e348c81c3e09c4dbfc4b0cf1e99166450d28e2d366ef7042b4a6`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.37.0/doc.go (source SHA-256: `ab02a40e5f32aec81867591d18e388f64f06ac238f3fbdf7ec33e2208ec1cbeb`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.37.0/error_type.go (source SHA-256: `9e0df41c1b73229305c3cf2899602574781be7a76faeea5f6e1687c3c05ea171`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.37.0/exception.go (source SHA-256: `d79d8ce58555b27816347fe2854a36f40604dea5423835baccd674904174fdcb`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.37.0/schema.go (source SHA-256: `031e428c9c671e233afc1b7c73e6b2b6cac5a1f247aa8c594224e899dbee12f0`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.41.0/attribute_group.go (source SHA-256: `2fb912eb9ef9730beb8490256adc4ff21b58aa8d03df05875e07d87f31a47954`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.41.0/doc.go (source SHA-256: `ab6f1c94472b932b962818aff7643221ba1e1fb9be1fdd30ed598ef9b447fa36`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.41.0/error_type.go (source SHA-256: `5ba647f889a16db8aa4bbd528315b5a94fa372c9c1b5fac703b87e5c549dda91`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.41.0/exception.go (source SHA-256: `c64aefe680661d47f9b6a92af15a68fcfd6e45e5db12d9d7f190f308e26e7f64`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.41.0/httpconv/metric.go (source SHA-256: `72be342087f3f403290020a67c90d9b8e96792578ae0e0db7a06c01cb1b6d03c`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.41.0/schema.go (source SHA-256: `069eca5ca954510660ab87d1bb694d170e550eb6041cdaa1cd76bcd3ba12fb67`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.43.0/attribute_group.go (source SHA-256: `ac0790dec17b81e57c772696cb79f696a5eef1fc6721f7fca62232736f3228c9`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.43.0/doc.go (source SHA-256: `2beae67530dd5ea51cf47b905196beccf3b2635678a1090bc07636cf67747e49`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.43.0/error_type.go (source SHA-256: `5ba647f889a16db8aa4bbd528315b5a94fa372c9c1b5fac703b87e5c549dda91`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.43.0/exception.go (source SHA-256: `c64aefe680661d47f9b6a92af15a68fcfd6e45e5db12d9d7f190f308e26e7f64`)
- go.opentelemetry.io/otel@v1.45.0/semconv/v1.43.0/schema.go (source SHA-256: `e0ab62c701a80ef2a68a9402ca751039b4476c6e203b1addcdde749a7ca50deb`)
- go.opentelemetry.io/otel@v1.45.0/trace.go (source SHA-256: `2d9a2025dfbfcd348d4ed3257ffd25dbe31e5e2385805c159b82facc6df1937e`)
- go.opentelemetry.io/otel@v1.45.0/version.go (source SHA-256: `0cf3c4496aec1f432f90ef5392a12d64b9e35422b1c63a2615c49246d95dd4e2`)

````text
// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
````


## go@go1.26.10

Source: https://go.googlesource.com/go/+/refs/tags/go1.26.10/

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

### go@go1.26.10/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### go@go1.26.10/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### go@go1.26.10/src/crypto/internal/boring/LICENSE

SHA-256: `56210f826b8f0fbac3160dfe55c97f4019eb6cdda3963b5a0eab8a8bdb62360e`

````text
The Go source code and supporting files in this directory
are covered by the usual Go license (see ../../../../LICENSE).

When building with GOEXPERIMENT=boringcrypto, the following applies.

The goboringcrypto_linux_amd64.syso object file is built
from BoringSSL source code by build/build.sh and is covered
by the BoringSSL license reproduced below and also at
https://boringssl.googlesource.com/boringssl/+/fips-20220613/LICENSE.

BoringSSL is a fork of OpenSSL. As such, large parts of it fall under OpenSSL
licensing. Files that are completely new have a Google copyright and an ISC
license. This license is reproduced at the bottom of this file.

Contributors to BoringSSL are required to follow the CLA rules for Chromium:
https://cla.developers.google.com/clas

Some files from Intel are under yet another license, which is also included
underneath.

The OpenSSL toolkit stays under a dual license, i.e. both the conditions of the
OpenSSL License and the original SSLeay license apply to the toolkit. See below
for the actual license texts. Actually both licenses are BSD-style Open Source
licenses. In case of any license issues related to OpenSSL please contact
openssl-core@openssl.org.

The following are Google-internal bug numbers where explicit permission from
some authors is recorded for use of their work. (This is purely for our own
record keeping.)
  27287199
  27287880
  27287883

  OpenSSL License
  ---------------

/* ====================================================================
 * Copyright (c) 1998-2011 The OpenSSL Project.  All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in
 *    the documentation and/or other materials provided with the
 *    distribution.
 *
 * 3. All advertising materials mentioning features or use of this
 *    software must display the following acknowledgment:
 *    "This product includes software developed by the OpenSSL Project
 *    for use in the OpenSSL Toolkit. (http://www.openssl.org/)"
 *
 * 4. The names "OpenSSL Toolkit" and "OpenSSL Project" must not be used to
 *    endorse or promote products derived from this software without
 *    prior written permission. For written permission, please contact
 *    openssl-core@openssl.org.
 *
 * 5. Products derived from this software may not be called "OpenSSL"
 *    nor may "OpenSSL" appear in their names without prior written
 *    permission of the OpenSSL Project.
 *
 * 6. Redistributions of any form whatsoever must retain the following
 *    acknowledgment:
 *    "This product includes software developed by the OpenSSL Project
 *    for use in the OpenSSL Toolkit (http://www.openssl.org/)"
 *
 * THIS SOFTWARE IS PROVIDED BY THE OpenSSL PROJECT ``AS IS'' AND ANY
 * EXPRESSED OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
 * PURPOSE ARE DISCLAIMED.  IN NO EVENT SHALL THE OpenSSL PROJECT OR
 * ITS CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
 * NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 * ====================================================================
 *
 * This product includes cryptographic software written by Eric Young
 * (eay@cryptsoft.com).  This product includes software written by Tim
 * Hudson (tjh@cryptsoft.com).
 *
 */

 Original SSLeay License
 -----------------------

/* Copyright (C) 1995-1998 Eric Young (eay@cryptsoft.com)
 * All rights reserved.
 *
 * This package is an SSL implementation written
 * by Eric Young (eay@cryptsoft.com).
 * The implementation was written so as to conform with Netscapes SSL.
 *
 * This library is free for commercial and non-commercial use as long as
 * the following conditions are aheared to.  The following conditions
 * apply to all code found in this distribution, be it the RC4, RSA,
 * lhash, DES, etc., code; not just the SSL code.  The SSL documentation
 * included with this distribution is covered by the same copyright terms
 * except that the holder is Tim Hudson (tjh@cryptsoft.com).
 *
 * Copyright remains Eric Young's, and as such any Copyright notices in
 * the code are not to be removed.
 * If this package is used in a product, Eric Young should be given attribution
 * as the author of the parts of the library used.
 * This can be in the form of a textual message at program startup or
 * in documentation (online or textual) provided with the package.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 * 3. All advertising materials mentioning features or use of this software
 *    must display the following acknowledgement:
 *    "This product includes cryptographic software written by
 *     Eric Young (eay@cryptsoft.com)"
 *    The word 'cryptographic' can be left out if the rouines from the library
 *    being used are not cryptographic related :-).
 * 4. If you include any Windows specific code (or a derivative thereof) from
 *    the apps directory (application code) you must include an acknowledgement:
 *    "This product includes software written by Tim Hudson (tjh@cryptsoft.com)"
 *
 * THIS SOFTWARE IS PROVIDED BY ERIC YOUNG ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 *
 * The licence and distribution terms for any publically available version or
 * derivative of this code cannot be changed.  i.e. this code cannot simply be
 * copied and put under another distribution licence
 * [including the GNU Public Licence.]
 */


ISC license used for completely new code in BoringSSL:

/* Copyright (c) 2015, Google Inc.
 *
 * Permission to use, copy, modify, and/or distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
 * SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
 * OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN
 * CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE. */


Some files from Intel carry the following license:

# Copyright (c) 2012, Intel Corporation
#
# All rights reserved.
#
# Redistribution and use in source and binary forms, with or without
# modification, are permitted provided that the following conditions are
# met:
#
# *  Redistributions of source code must retain the above copyright
#    notice, this list of conditions and the following disclaimer.
#
# *  Redistributions in binary form must reproduce the above copyright
#    notice, this list of conditions and the following disclaimer in the
#    documentation and/or other materials provided with the
#    distribution.
#
# *  Neither the name of the Intel Corporation nor the names of its
#    contributors may be used to endorse or promote products derived from
#    this software without specific prior written permission.
#
#
# THIS SOFTWARE IS PROVIDED BY INTEL CORPORATION ""AS IS"" AND ANY
# EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
# IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
# PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL INTEL CORPORATION OR
# CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
# EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
# PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
# PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
# LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
# NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
# SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### go@go1.26.10/src/vendor/golang.org/x/crypto/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### go@go1.26.10/src/vendor/golang.org/x/crypto/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### go@go1.26.10/src/vendor/golang.org/x/net/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### go@go1.26.10/src/vendor/golang.org/x/net/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### go@go1.26.10/src/vendor/golang.org/x/sys/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### go@go1.26.10/src/vendor/golang.org/x/sys/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### go@go1.26.10/src/vendor/golang.org/x/text/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### go@go1.26.10/src/vendor/golang.org/x/text/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/edwards25519/scalar.go (source SHA-256: `b9bc9ba5a2160dccee62ccbb4371c90c5d1c654d55b8ae08f33fa29c58876050`)

````text
// Copyright (c) 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/edwards25519/edwards25519.go (source SHA-256: `6bc3d5ea48276f302a7d55999f5c24986012f8d65ed1749b483f641bbef7bc18`)
- go@go1.26.10/src/crypto/internal/fips140/edwards25519/field/fe.go (source SHA-256: `7e678ceaf12dff6b4f09ff09e45d02fe912f55d8147c527ba5865d2c0007aea3`)
- go@go1.26.10/src/crypto/internal/fips140/edwards25519/field/fe_generic.go (source SHA-256: `75d213650f76b7c5501202bf5cae1f0cdc1f84fbaa8bed5f93c0109189e01dba`)

````text
// Copyright (c) 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/edwards25519/field/fe_amd64_noasm.go (source SHA-256: `3f3d670d4da11ac2e68a8dfd69d869f955ea9a85264cf88ee55184e61f74ba2e`)
- go@go1.26.10/src/crypto/internal/fips140/edwards25519/scalarmult.go (source SHA-256: `11f525c854f8e9738a5779adf61d8b3cb39067bd9ccad6a908355811242ecc1f`)
- go@go1.26.10/src/crypto/internal/fips140/edwards25519/tables.go (source SHA-256: `fa11b18ed84768605ce4b532367295199a32a211a070f5a65b38048f298d653f`)

````text
// Copyright (c) 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/edwards25519/doc.go (source SHA-256: `f4fb8e9a475e0e98295684542301dac6087bf799da02011f920a567ecdc63d8c`)

````text
// Copyright (c) 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/bufio/bufio.go (source SHA-256: `0393162eb1732d74ec5efe88760df58186f6847c525e577b297e8649c679580b`)
- go@go1.26.10/src/bytes/buffer.go (source SHA-256: `4457385c174f963728d57902c5b474dacf5c2b7226f103045371a2b4cde7f119`)
- go@go1.26.10/src/bytes/bytes.go (source SHA-256: `2228844174e0f9b5e751bade463cb9ad5f0cd8ef209f5cabcb3cc18d8844438c`)
- go@go1.26.10/src/compress/flate/deflate.go (source SHA-256: `4da50a0ab834a5e152ae337f9fea35f3527e29ea415e629dfe84749dfec07388`)
- go@go1.26.10/src/compress/flate/huffman_bit_writer.go (source SHA-256: `ea6ea80d176d172a45b072190a66237a86c10969411514b07feff614fd05cdce`)
- go@go1.26.10/src/compress/flate/huffman_code.go (source SHA-256: `84b830cd8cc3da67f7572d7d8f244f6be0ea21ebd2f579753e91ab29870f5534`)
- go@go1.26.10/src/compress/flate/inflate.go (source SHA-256: `28b241ba6b66cc14ac93cc17a49756cd78daaade60afb47d7cd6aed2489dc5cb`)
- go@go1.26.10/src/compress/flate/token.go (source SHA-256: `bc96536b4ca6e946a5094c3c60a085485cf6e25be1fc74914bdefcf4352f1516`)
- go@go1.26.10/src/compress/gzip/gunzip.go (source SHA-256: `5886344b7fd61358d430872ab735653b16a314c87dab1f0bde4d8d9db0381ec4`)
- go@go1.26.10/src/compress/zlib/reader.go (source SHA-256: `4c669325ef280fa8dba1878b667577bb0229c0f1b91b98c2dea9165b9430254c`)
- go@go1.26.10/src/compress/zlib/writer.go (source SHA-256: `24b0f31cc40b4937452c64e80a03da246ef078967070d7d84f1660d2daa5a488`)
- go@go1.26.10/src/container/list/list.go (source SHA-256: `88d1eadd6ac199fe42872bf52ff71264235fa51b0c5034be3d099fcfddcf640e`)
- go@go1.26.10/src/crypto/aes/aes.go (source SHA-256: `582ccec90f9838939c00400941e427f7c305362bd96d91528a64b68bafaf72e6`)
- go@go1.26.10/src/crypto/cipher/cbc.go (source SHA-256: `044961bab9631d5c9a558dbc6e4fd4199c362d60f4887599b119ab296388c9c7`)
- go@go1.26.10/src/crypto/cipher/ctr.go (source SHA-256: `eab2ecce91dc7ba343ba116709c5656e5cab61bbf96fb359b3908bfd83343fb0`)
- go@go1.26.10/src/crypto/hmac/hmac.go (source SHA-256: `91ce21eb1c6889f5e30ac98fcf2846d5c75e76cf3aa9ea9a74f9724a49b10d20`)
- go@go1.26.10/src/crypto/internal/fips140/aes/aes.go (source SHA-256: `3f7a470932a6f1e2bc9fa5f03be3a60f7c2edd557ec4a58ecb25fcab71b83ddc`)
- go@go1.26.10/src/crypto/internal/fips140/aes/aes_generic.go (source SHA-256: `a33d0a61c315c185f168987c7a2f4035d78b1059e9981cf750b7bdf15f91af08`)
- go@go1.26.10/src/crypto/internal/fips140/aes/const.go (source SHA-256: `1e425707426f310a093c92512a88c134f508801dae0326b1deea4c7e135fdf38`)
- go@go1.26.10/src/crypto/internal/fips140/hmac/hmac.go (source SHA-256: `50de5ed25e3b168be1853876a6c8170ea42079e43ce1a4fc1a35a60e51c14234`)
- go@go1.26.10/src/crypto/internal/fips140/rsa/pkcs1v15.go (source SHA-256: `74af6dcc3e178f71d51ce5fe624df8d6bf7c5e1e36bca92096f65948ee197935`)
- go@go1.26.10/src/crypto/internal/fips140/sha256/sha256.go (source SHA-256: `94f8a5508cda46bc036b891d32eae19eb31253946d3442e44d4bb3e0341643b8`)
- go@go1.26.10/src/crypto/internal/fips140/sha256/sha256block.go (source SHA-256: `eb49287421c308743473db444fff531f65231f1c5d8ecbc48b6b4ca1441975ec`)
- go@go1.26.10/src/crypto/internal/fips140/sha512/sha512.go (source SHA-256: `2d5b0c8c656c71c5a9eff37daf1919141e7742f519c7ad135336ee8805f76cdc`)
- go@go1.26.10/src/crypto/internal/fips140/sha512/sha512block.go (source SHA-256: `4c5febb80a02435b3b9baa16b6009dbb7ed4dbc6027455bbe7eb609109803ffb`)
- go@go1.26.10/src/crypto/internal/fips140/subtle/constant_time.go (source SHA-256: `614785e04c5e9ea725b90e1584bd79011714a27458f309bb48c6a60a6c68513f`)
- go@go1.26.10/src/crypto/md5/md5.go (source SHA-256: `04d072c18c5782f8f94f3687e731c0a2c53d8bb695a3b1c334425067ac181c17`)
- go@go1.26.10/src/crypto/rc4/rc4.go (source SHA-256: `57d37b459eb60300b0e7979f04c9cf680a15f5da87a5941fd3e0e39302c59d76`)
- go@go1.26.10/src/crypto/rsa/pkcs1v15.go (source SHA-256: `ca0ccb3ea55a56ddb169ab3874d348c349bdc8b40ce8afa9bb6d7189dc8f78b0`)
- go@go1.26.10/src/crypto/rsa/rsa.go (source SHA-256: `648123c00aa7c017facd2e5a8d5c207b6fdd3534a088bda64307a633809c14f0`)
- go@go1.26.10/src/crypto/sha1/sha1.go (source SHA-256: `ee286c960d9f42c93e95c59838e74d4420641497fc3eb6e84ac2cc2d2e367ae7`)
- go@go1.26.10/src/crypto/sha1/sha1block.go (source SHA-256: `64a5f9acf7a3dc268abbcd16e6727a8e105267393cd909c2736ec8d8f32a24f3`)
- go@go1.26.10/src/crypto/sha256/sha256.go (source SHA-256: `34c5c7e8408022e7b8fd570453aebe6a95607761c14fbae61f465f9ba6328da9`)
- go@go1.26.10/src/crypto/sha512/sha512.go (source SHA-256: `f2edc5fa3205fd4e2bc39d0b72a5dc368d8e5b7ab65816fe43e7db62ab495c50`)
- go@go1.26.10/src/crypto/subtle/constant_time.go (source SHA-256: `34a907482c81c783f55fbbeddb8099859e9d6c304c623f58d1f285fce4476940`)
- go@go1.26.10/src/crypto/tls/alert.go (source SHA-256: `95b0aab666e8a08a2c56191e0ed96a804be83b4e2d7475d4d1e7c4eb9f8b6b6c`)
- go@go1.26.10/src/crypto/tls/common.go (source SHA-256: `d2837fbe55a398c7362b3ff8ffe43c06d0832df2305f45f32aeaebadc784486d`)
- go@go1.26.10/src/crypto/tls/handshake_client.go (source SHA-256: `7d6210c69ff9bf0d8506a8f2f59bf33a6132d0a3574487bcc356f858e50c6fab`)
- go@go1.26.10/src/crypto/tls/handshake_messages.go (source SHA-256: `b3a61247d78d93e196debfb6d6202e9f771d732de94bd4c30d7dbbebad4ac2b3`)
- go@go1.26.10/src/crypto/tls/handshake_server.go (source SHA-256: `371e52320f1011df2b23f1ee942f0ef0e6130cb53c8e7b2f83b7aff69cae21cb`)
- go@go1.26.10/src/crypto/tls/prf.go (source SHA-256: `2718f4cbf933bf5488396852fe0ec7795daec546e3708e43adf551f8b4c322d5`)
- go@go1.26.10/src/crypto/tls/tls.go (source SHA-256: `f7496eda12403ed3008aac8f1bdc7aa0238e4f9898eb5fe7bbbc1af711203554`)
- go@go1.26.10/src/crypto/x509/x509.go (source SHA-256: `6e9fe1d0355d04cfe9def570906a753b9036d1f251b714b3901bb24d24e077ce`)
- go@go1.26.10/src/encoding/asn1/asn1.go (source SHA-256: `325e3acb01e6ca0083bfe78cdd720aeca0d766c8561c4125e8e1711358d25740`)
- go@go1.26.10/src/encoding/asn1/common.go (source SHA-256: `302a3808ff7c006e0fc8d0342b0ee4196d08ef9f43a3b77268345f20e7e2aa15`)
- go@go1.26.10/src/encoding/asn1/marshal.go (source SHA-256: `32e7bc36f97585c359794908f7cf7261b1c6e6a0515b8a2f5c961b8991670b9b`)
- go@go1.26.10/src/encoding/base64/base64.go (source SHA-256: `e6437c7716ecd1ae91c9d10037362cab186e5873a6365b713e2d8c57d8fc32fc`)
- go@go1.26.10/src/encoding/binary/binary.go (source SHA-256: `fc4fab16ee184c5eaafb75c2e2b8728c9dd04228f96ec8d394c9d3110f9c91df`)
- go@go1.26.10/src/encoding/gob/decode.go (source SHA-256: `99c8b170d943f63d0c7a63b06774cb574238faf76ddca738d81a91ccc89251f7`)
- go@go1.26.10/src/encoding/gob/decoder.go (source SHA-256: `e6bd8abc697c5db332fc69d99aef5e973db4074b63605f25bd4bcda18358c633`)
- go@go1.26.10/src/encoding/gob/doc.go (source SHA-256: `edda0f10e817f4a6fa5014ec131940fce9007d4b01b7f3fe6328e5db8ad3b297`)
- go@go1.26.10/src/encoding/gob/encode.go (source SHA-256: `82d3bf376aa83c0401e1086891478764cba9bf4eae00b7cc62b7aecbcaf11d6b`)
- go@go1.26.10/src/encoding/gob/encoder.go (source SHA-256: `d9a5c9529ebe5d8891970427ffc8049c4ba9a387d12c2ac800f0af7a376366c1`)
- go@go1.26.10/src/encoding/gob/error.go (source SHA-256: `b4405a850719ac039c2c5b56ebec91421c18cffc6f4d1aba655d85f944c98998`)
- go@go1.26.10/src/encoding/gob/type.go (source SHA-256: `e151a4d8ac10339dd0ab2c6c40edb081b8ea8b0fb3fe4c2581f68b54624baf04`)
- go@go1.26.10/src/encoding/hex/hex.go (source SHA-256: `b3344957ddc5d4090972f94b552a5a7a5bdfb8174a189df033adb0ede97002ba`)
- go@go1.26.10/src/encoding/pem/pem.go (source SHA-256: `536954f803f79d8972e8f86b792b25d4ac83167fbe4c3117954a4378e60521da`)
- go@go1.26.10/src/encoding/xml/read.go (source SHA-256: `d0e759984eb7e8a48cdbc72808a81a9ff59a6e92c8a5186aae9a7db4d105357f`)
- go@go1.26.10/src/encoding/xml/xml.go (source SHA-256: `41a7f2c7c910a57f70ce460cbe057a31703e11dbb169b499e50a3339a5f33bdb`)
- go@go1.26.10/src/flag/flag.go (source SHA-256: `eb3238110346601e8e014dac5e5fa873bde7c94bb2f33053f92fd175c1b2a1f1`)
- go@go1.26.10/src/fmt/doc.go (source SHA-256: `400732eb804a960f2c1bdfbd88a5cb183e11e7a46bbb44f2b0784d0b1b171a05`)
- go@go1.26.10/src/fmt/format.go (source SHA-256: `a668c0d6bb23f6ac929b73ef661a65be2225a9f665052a2ce02c4c58339c3335`)
- go@go1.26.10/src/fmt/print.go (source SHA-256: `57250462a0d06d61f124b4530493bcdd150e9de29600bd60347e1c74457ac8a1`)
- go@go1.26.10/src/go/token/token.go (source SHA-256: `fac9e32b0a22562278d9fd23faada551d7c9f05dff2787d6e7dcec3f7ca18a23`)
- go@go1.26.10/src/hash/adler32/adler32.go (source SHA-256: `9301031a6cadbfde857a96d39f3b9c7d8455aec45f486592627c04a514c75443`)
- go@go1.26.10/src/hash/crc32/crc32.go (source SHA-256: `85f9b1484d31d449fda4c70b6df4e8c78455b5384d06fecd5f7a359638de2f76`)
- go@go1.26.10/src/hash/hash.go (source SHA-256: `2784a10b0be2f58c2a0157a5457df5a8e3ad4dfebce5dfb118b7df22ab812776`)
- go@go1.26.10/src/image/draw/draw.go (source SHA-256: `24f42de1373014deefce87d4dd02856010ac8d385c4e61eb0c52a1a1273b0207`)
- go@go1.26.10/src/image/image.go (source SHA-256: `ec22e85f80321c25d645e046d6378a919b0f3d568c2af101226c1c93b953f39b`)
- go@go1.26.10/src/image/jpeg/huffman.go (source SHA-256: `0dc31ab00dcdfc9853d11ff56cce8ce5ff2099e8b2ed18ab35561aa12993897c`)
- go@go1.26.10/src/image/jpeg/reader.go (source SHA-256: `9af311c6149311b0863a3f3175ac7b5ecb7f86b487006a6ba65f0aa69f7ddb3d`)
- go@go1.26.10/src/image/png/reader.go (source SHA-256: `f6e97084ce561f92ae7658b34f16d797aa46c65325b2ae9cd931d0f2434993c1`)
- go@go1.26.10/src/image/png/writer.go (source SHA-256: `88f6b695cea6fb90f8b96fd837a4f777fe7d2d7c97d88741e3a6fa39c3544784`)
- go@go1.26.10/src/internal/poll/fd_posix.go (source SHA-256: `ff002941e536f405d2876006c2973a1b61d26fd032094ade5ed2fc289fc35084`)
- go@go1.26.10/src/internal/poll/sockopt.go (source SHA-256: `ee7f52272fdecfd2e40ad6eb0616b49c781c39c1393dd0c47a11e6bc2caec1b0`)
- go@go1.26.10/src/internal/poll/sockopt_windows.go (source SHA-256: `9935127f68398d47b73f7eadddfae21caba27dbaff99581e4a5e4283e3d7e2a9`)
- go@go1.26.10/src/internal/reflectlite/type.go (source SHA-256: `5ce048cf76070f0c0734694bc5c15e96886bd44f726c32f5f72ca33a1c5f94dd`)
- go@go1.26.10/src/internal/reflectlite/value.go (source SHA-256: `37526f8c1c775bf455060cd93c6cfcc2518ba0de38f63d78cdd1def75b597212`)
- go@go1.26.10/src/internal/runtime/atomic/atomic_amd64.go (source SHA-256: `2b651194b7a1c268dbbc092185afd377d49cbc9c66208d25b30b9e936e7b0bbd`)
- go@go1.26.10/src/internal/strconv/atob.go (source SHA-256: `6938f9a33a3daf57d7bada9b6c38b4cc26487c3e613c42c3512a055bed10d6f9`)
- go@go1.26.10/src/internal/strconv/atof.go (source SHA-256: `2c2205ccab8626b623a188144a7dcaf202ece77002a27a2cd6b0aaf5a61669a5`)
- go@go1.26.10/src/internal/strconv/atoi.go (source SHA-256: `eb5ad640903abe9e8c8b7fbce8dd3487176dbeb43e8b8b4efa6e73ecb5207ddc`)
- go@go1.26.10/src/internal/strconv/decimal.go (source SHA-256: `23bd45989d0b5a7b441c45873bf43f210fdf8381132fc74424d381216ac3a8b4`)
- go@go1.26.10/src/internal/strconv/ftoa.go (source SHA-256: `2c12d1c3a17b9a978266a7c761862a4be8dd1f2fae14ad5bbdf6ff2dd8f3493d`)
- go@go1.26.10/src/internal/strconv/itoa.go (source SHA-256: `cfda0d39152dc8ad8e2eb0407682c399ea0a321796211e015096fa57b17a35d5`)
- go@go1.26.10/src/io/io.go (source SHA-256: `3945d328a9072e4f5970b7d4bc0eb834eec29fd22889a1ed9f3fb76ad41c3274`)
- go@go1.26.10/src/io/ioutil/ioutil.go (source SHA-256: `166c114baec89fd662376a1f08682589e97f595023df450a9f17a42ae419af39`)
- go@go1.26.10/src/io/pipe.go (source SHA-256: `75ef56f505592fcffb188bf7660c6ea6abede3f2701b500d6946fa519c37e233`)
- go@go1.26.10/src/log/log.go (source SHA-256: `57ac96bed54d29c3ed472f6af3001a6258021d123363068d3630875b7704ecf0`)
- go@go1.26.10/src/math/abs.go (source SHA-256: `9608d0f0e20d941025c2bf268e8541dbecaf7e9b1524e030c390166384a7da96`)
- go@go1.26.10/src/math/asin.go (source SHA-256: `846e55289485807288cf67f40a20c7f16c246dd831a784fec3e88e7af9c96a4c`)
- go@go1.26.10/src/math/atan.go (source SHA-256: `777b61ebf72d42b754004fe6a5e9a4ac7fb85fb95cd5e6fef396540ce6c42ee1`)
- go@go1.26.10/src/math/atan2.go (source SHA-256: `5a928623ed08c894f79463e5941671cf889b925c1d44eb69842f4cb5fa65406e`)
- go@go1.26.10/src/math/big/arith.go (source SHA-256: `de8c0f132a42d2cabe62ff14e42eca845fb2fd6b27f6db9b42afc2ad7aefff87`)
- go@go1.26.10/src/math/big/doc.go (source SHA-256: `cc7e6416479405ff79a09636d812e78be84e1ee5b225df217a34c69b709a7ce5`)
- go@go1.26.10/src/math/big/int.go (source SHA-256: `7689a5447c12d8c4cb350697cd7cdb4c827b3bdd3f0fb875441588cde3b18676`)
- go@go1.26.10/src/math/big/nat.go (source SHA-256: `9d2cffd653f28a983014397ad220867c49a0a31f14aca30ff0a653638c15b955`)
- go@go1.26.10/src/math/big/natdiv.go (source SHA-256: `c2400503fd20c78a98b41a6c7aca1307a67b5031e8c5459a6e2395b2f7da2a2b`)
- go@go1.26.10/src/math/big/natmul.go (source SHA-256: `f1767186e544a77fea89923a0051db1f3a36658d34f12ef4b35b78bb42234d40`)
- go@go1.26.10/src/math/bits.go (source SHA-256: `eaf5a27831d55da48238beb6a86b08e8ce209cd36a535a24a457f7546ca031f7`)
- go@go1.26.10/src/math/cbrt.go (source SHA-256: `de08f8225477bd143bd5318c5f46b7e26bdb008db8e4d174c6130e751e66e335`)
- go@go1.26.10/src/math/const.go (source SHA-256: `859fedd369b0359b356ee4c54345a8d528e8a19f69c2eb93e7256a073689fe5c`)
- go@go1.26.10/src/math/exp.go (source SHA-256: `75d415ac8eca1afc796cd621b393153ff82d342c24f5b391dfdd6a6e70303f2b`)
- go@go1.26.10/src/math/floor.go (source SHA-256: `86315519ca54caa22134e79c02126704a67c00f68dd25ea8289e958278060033`)
- go@go1.26.10/src/math/frexp.go (source SHA-256: `5f11ebcb154799d67e2d15c13673c8812557994be6141500edf0358550cffab7`)
- go@go1.26.10/src/math/ldexp.go (source SHA-256: `d928bbad23be48527ee0c7255fd9ce5f35b60dc958d60289b5fcdfa8066524b4`)
- go@go1.26.10/src/math/log.go (source SHA-256: `d8675d3d9c665ec0db5cc15478c03efb244301bfa1268d02ae1861bf6f2812ab`)
- go@go1.26.10/src/math/log10.go (source SHA-256: `936bbdfbe82a5b8028bd55b92666792954f3b8908a8e4e27008ec23935c96e37`)
- go@go1.26.10/src/math/modf.go (source SHA-256: `03f1cdef8c6a368d60a88a90edf5e0b1720f3789a068a6b74b840c28c9bb037d`)
- go@go1.26.10/src/math/pow.go (source SHA-256: `848a205994a0a507046a91a02ad4a95172af2fa5f6f11847acb8cb634471b670`)
- go@go1.26.10/src/math/pow10.go (source SHA-256: `75c48d71fbe5930b75f736688145efa2f6373938ce988b7bd659d5ee273378d3`)
- go@go1.26.10/src/math/rand/exp.go (source SHA-256: `e3837671076ec0432cf4347c4154f87c88c15b2e2b6a739690a68e2749807889`)
- go@go1.26.10/src/math/rand/normal.go (source SHA-256: `d8f4450a14b03bd0ec939eb4b127e199d9d78ba3d07c4ad0a65ac493dfa1c62a`)
- go@go1.26.10/src/math/rand/rand.go (source SHA-256: `a57f58bcc4baa4ac71fb09a6f2a8ce525f6f097d7b9c7ea61f206bcfd910aa39`)
- go@go1.26.10/src/math/rand/rng.go (source SHA-256: `e301d938fe5dbe4ef5ff52f31de68d2547c3ab50da5c95bb21c49ea983a817ab`)
- go@go1.26.10/src/math/rand/v2/exp.go (source SHA-256: `1962b9394fb390cdbbba401ca7bb6c0e10b8597d095890ec6421811d18ff2f23`)
- go@go1.26.10/src/math/rand/v2/normal.go (source SHA-256: `355b9570ff06ec4825b4fb8a447e71208ac38ecacde0c289f88bf524fd85fa89`)
- go@go1.26.10/src/math/rand/v2/rand.go (source SHA-256: `a0b0f516652e66ee56de089ec01457f5f786e0463cd08b2453eb025cd6ad6f58`)
- go@go1.26.10/src/math/rand/v2/zipf.go (source SHA-256: `ee4616e17c8b1d7ff33229b075f2e73c369307be02258f173393508cb2db0357`)
- go@go1.26.10/src/math/rand/zipf.go (source SHA-256: `87f102ccb820854fb6b7e2b86f78fe793de8fd79adc2bd316c082daf70267326`)
- go@go1.26.10/src/math/sinh.go (source SHA-256: `325a233cb3a9f01b214a82a47ad438795122cf88f2f9e04e9fdf6819f394acda`)
- go@go1.26.10/src/math/sqrt.go (source SHA-256: `c7e036d1511bea94992014f9b8c9259e5737c385cdd5135db9092369250d4a94`)
- go@go1.26.10/src/math/tanh.go (source SHA-256: `a80f128fa85ddc9e94983b53f4a928f6811b96174e7920032e3e1a4367b96591`)
- go@go1.26.10/src/math/unsafe.go (source SHA-256: `dfe6fe17844c02c370d06d7d7ca824ce167f8dc010ef79b692f4a14fa1ffa896`)
- go@go1.26.10/src/net/dnsclient.go (source SHA-256: `a73f909dbd747cd43eb1118d820cc6674f7147a263d02683281c658db033f3d3`)
- go@go1.26.10/src/net/dnsclient_unix.go (source SHA-256: `97e1643224272588e5161f0572d3a36d80c4a636c1f6226ce7bc5b610144b630`)
- go@go1.26.10/src/net/dnsconfig_unix.go (source SHA-256: `a707636ba4d5c9fab72cf500a0bc9d617385fee6fa21d6a85bf5ef41fb250748`)
- go@go1.26.10/src/net/fd_unix.go (source SHA-256: `0f8f72e4a3bafa38f3a12a826f1ffdaa4797c921a36e152d4e991e222ea8d3e0`)
- go@go1.26.10/src/net/hosts.go (source SHA-256: `265d5f3d71175407a3cc49c731d43ded8e90fdb2a67106d3d4ecb0e86955c3eb`)
- go@go1.26.10/src/net/http/client.go (source SHA-256: `908d9d964cd157899b60e7ba839befa4bfd398b4fb8f7ff5a1aa389eae75f652`)
- go@go1.26.10/src/net/http/cookie.go (source SHA-256: `a0b5b5a3be97de79e68a6bcd9df2c34dfb9bd09c6eb312b38ca1bbc79cca86dd`)
- go@go1.26.10/src/net/http/fs.go (source SHA-256: `1cb8e594e15c1694f5f56dfab0ebef55b862bf1de599ab6fa9d216e5a1770126`)
- go@go1.26.10/src/net/http/httputil/dump.go (source SHA-256: `6e735a669ee4168fb6d27b68bcbe684af88f714e87973c71614e18484e31dd2b`)
- go@go1.26.10/src/net/http/httputil/persist.go (source SHA-256: `1a93590cfef8b9e6211ce6addb5145364fa4ca67e1390480b3096faf9f137389`)
- go@go1.26.10/src/net/http/internal/chunked.go (source SHA-256: `d2a088e8d1d0ae6ff593b27011982ef9115fee593cb57cc97eb34982ff17f5e2`)
- go@go1.26.10/src/net/http/request.go (source SHA-256: `4a34a9e157a606a046d5b9f8f761c571c4762668fa351d92048515a595939c38`)
- go@go1.26.10/src/net/http/response.go (source SHA-256: `43833d3f5b76380a6ba2ac6972d024f9efc599fbd17ceebda34e3e79b03483ec`)
- go@go1.26.10/src/net/http/servemux121.go (source SHA-256: `e984a617731e8822da83f9619b2063a264c84aef9ab81bfad13665909eff4536`)
- go@go1.26.10/src/net/http/server.go (source SHA-256: `6f344ca40661ce5d31ce7f237dc23c53d4d34fe3043bd5e76bac50e6c97d91cb`)
- go@go1.26.10/src/net/http/status.go (source SHA-256: `7633ca0ca545767f3db3d8960a96ee4ecd68fb39439c88b7c03efcb3d1786297`)
- go@go1.26.10/src/net/http/transfer.go (source SHA-256: `3e04744176d3cf5bef2827ab158b10aaeec3981ac5f4e7fb86e3dbf5f592eed7`)
- go@go1.26.10/src/net/ip.go (source SHA-256: `7e26071410dd7265f58c27e9a13820b9e42b13187529db3f4a43eb48f94d1436`)
- go@go1.26.10/src/net/ipsock.go (source SHA-256: `1d5cc127665695ab1351a96b015ceeddcc5adcbed91f35af0a6c6cd1552d846b`)
- go@go1.26.10/src/net/ipsock_posix.go (source SHA-256: `9d2b4164bff19dcf053ccf2298790b9d96c604c4c374bdb965f367c3d3aa82f0`)
- go@go1.26.10/src/net/lookup_windows.go (source SHA-256: `9c7a89fbfa9241a64e4cca39d0155b19e708efb6eb1b996561938bc71cff50a8`)
- go@go1.26.10/src/net/net.go (source SHA-256: `35a3d5047c156a1395c4c197fad5a99e123567a12ae7bc08ba7d5c723281ae03`)
- go@go1.26.10/src/net/parse.go (source SHA-256: `5ae9c14a1300f39fb0de55a098af5f0b679b4ae8f1677b86be35e78264378cda`)
- go@go1.26.10/src/net/port_unix.go (source SHA-256: `34613ab7f4204a7a08cfc58e12b612f6ab544a46ad55484a687c34419f3f5576`)
- go@go1.26.10/src/net/rpc/client.go (source SHA-256: `e2c761997a75c3bf592b5a1d70252633f7f8a44c615307e194ecb8f727a7ad08`)
- go@go1.26.10/src/net/rpc/debug.go (source SHA-256: `7bf7c8ea236796a3ad01745b99a58769bdff9cd981c26b1a85535a1d88b75ca8`)
- go@go1.26.10/src/net/rpc/server.go (source SHA-256: `b10f3723b6d40543d3feed65f0a085258114fcc95c119023cd611b178942e7a6`)
- go@go1.26.10/src/net/sock_bsd.go (source SHA-256: `f56d294f7cceca2478bd066fbb1df817e27678dac53d1d21ea30873dc5e8259e`)
- go@go1.26.10/src/net/sock_linux.go (source SHA-256: `ab937db35b43c476ca474fb4fbe6a929547485f30488a16bc9a0cd26e4404888`)
- go@go1.26.10/src/net/sock_posix.go (source SHA-256: `f83091e2fe19ea25d190115531b68f76f128c88d817f6c9d294ffbc7fb154083`)
- go@go1.26.10/src/net/sock_windows.go (source SHA-256: `c24d90f9a90da5aefbe40d45c84c26e1ea5869a20e142f1f1101e7fd0f2f981d`)
- go@go1.26.10/src/net/sockopt_posix.go (source SHA-256: `eedba1ea662b37dde58015af6253b84702840a5c9ef5e25a20485525e94a52cb`)
- go@go1.26.10/src/net/tcpsock.go (source SHA-256: `0c3392bd7c2966f894e01ea27514623cf09dc2660e426049219f7e5ee3ec0f74`)
- go@go1.26.10/src/net/tcpsock_posix.go (source SHA-256: `01136ba949dd7d7f70511d514ec12a89a22c216457b2ab6220c16387bb62dbed`)
- go@go1.26.10/src/net/tcpsockopt_darwin.go (source SHA-256: `9a38de57d27c95009e59f218811ff3aee85cdc812134ef2ab1e44a810db71b25`)
- go@go1.26.10/src/net/tcpsockopt_posix.go (source SHA-256: `16d8565b480a743e5d75ccc457fbb2dc46f4fcf3a681ea01bf45be0f057e348d`)
- go@go1.26.10/src/net/tcpsockopt_unix.go (source SHA-256: `078d4f4b93ddb98e0f668e8652de2975eb768499d0a8c3e064f0d7e3c66f50ff`)
- go@go1.26.10/src/net/tcpsockopt_windows.go (source SHA-256: `4492fd3fe6d4e7e131a04547f0f2e6c7ea36bae8d442569112e1168fed401c5e`)
- go@go1.26.10/src/net/udpsock.go (source SHA-256: `1a05ced29d6c7eb82b12b480a37739f67e695be45fa340f9c5bc908b14f5bf1a`)
- go@go1.26.10/src/net/udpsock_posix.go (source SHA-256: `374963d9933fc3447905a2958cd0179a3559ea0b4016bb05deac5b26dff629c3`)
- go@go1.26.10/src/net/unixsock.go (source SHA-256: `88ab17cda1f24aba697dba9f7ca86d9becfdf6f5a0bef808a9d0c42ad5ce5ed7`)
- go@go1.26.10/src/net/unixsock_posix.go (source SHA-256: `472016972215bcb170c588bab5d36fb920e071e64c0b595d68da91a45147d561`)
- go@go1.26.10/src/net/url/url.go (source SHA-256: `7bec001746c2cd512275c7269bc293f3d36cb004ddb193d736e02e361a3778aa`)
- go@go1.26.10/src/os/dir_darwin.go (source SHA-256: `bc7c50b0d6452d6b5fd1ed0b06e8b8341a476e5d11af56747717f4233cd58c67`)
- go@go1.26.10/src/os/dir_unix.go (source SHA-256: `0de1bd5b40affb85fb104b4f04589fcaad093cd4b8165f6b2161427f7935bad0`)
- go@go1.26.10/src/os/dir_windows.go (source SHA-256: `bed7533115f62133558144ff8e3e5944a8c6f02338b4fbd29f453d85f6145a8c`)
- go@go1.26.10/src/os/error.go (source SHA-256: `f073f7fb0bad2cd69df450956b4a267b5ecfb519a9089d44fe6b6e2d6f766e18`)
- go@go1.26.10/src/os/exec.go (source SHA-256: `b41760eafad2ed3e1cf68ad75823c222ba2f78d683cc25e81fb8eb4f96873d2f`)
- go@go1.26.10/src/os/exec/exec.go (source SHA-256: `0bbfd988045f66e78e7416a1a0d47016ea966c65323d2c6ca7a18cdf4b7b5df0`)
- go@go1.26.10/src/os/exec_posix.go (source SHA-256: `57f9028ac13e80d2106dc60649b0101b540957729a1dc9f775211aea6bd9a624`)
- go@go1.26.10/src/os/exec_unix.go (source SHA-256: `6e5581c9945c2364c01c26ed4d8659f0b0c10cd32b4eab3000132dc98e8e92f4`)
- go@go1.26.10/src/os/exec_windows.go (source SHA-256: `7a3c8f10605e5ea468538b1833b04ddf609bff84999b037839ef23d99afc50c4`)
- go@go1.26.10/src/os/file.go (source SHA-256: `85a4adcd383c9947d6133b65ccd01e95b68ce1d07b8474b7b7aad8ca9f3b8139`)
- go@go1.26.10/src/os/file_posix.go (source SHA-256: `2e93ac5f812cd7184d3ac3a2eb9a51583911a6261cf9c27772fa7dfedcbe514f`)
- go@go1.26.10/src/os/file_unix.go (source SHA-256: `877dd004355c475e8805b5700061604558422da057037556e0cf6edeb36ab7ec`)
- go@go1.26.10/src/os/file_windows.go (source SHA-256: `c22ecd7dd18b10a3091f7c1d7deae868a64e7adcf046c8ca62d3666e73edfcc7`)
- go@go1.26.10/src/os/getwd.go (source SHA-256: `20bf5c0875341ff1d685ff8aa58f3a0de19250ec19708d32a921779ed515e26e`)
- go@go1.26.10/src/os/path.go (source SHA-256: `a901784cb6abfb4817a379bdc909f42090b9cc5e49a6b2dac5c06ee6e844966a`)
- go@go1.26.10/src/os/pipe_unix.go (source SHA-256: `afbd7c5f5d714feb8b72739913311719e3c352ad2c21b70b9c043106da8a9bb1`)
- go@go1.26.10/src/os/proc.go (source SHA-256: `54ed65b25fe5251adb4b7799117f38ba0b2176610f2d21b333a7a5745e988200`)
- go@go1.26.10/src/os/stat_darwin.go (source SHA-256: `2e2822d6b5ff3bab732488f16ba2e2e4ff32f214efa9319b86ab3c3bfff0fab7`)
- go@go1.26.10/src/os/stat_linux.go (source SHA-256: `62cb59a3767a69fbb717c96f6993d1e80eeb0f1660fc452bfedd01100339c7dc`)
- go@go1.26.10/src/os/stat_windows.go (source SHA-256: `90f01ab7b8a973df1c4588c3b40fae737c6d23d6548a1805e5f213e38def6df4`)
- go@go1.26.10/src/os/sys_bsd.go (source SHA-256: `0a5cd4a328bb028c01b1e41a9f891b9176229234c6f4e18afeab9fe771392193`)
- go@go1.26.10/src/os/sys_linux.go (source SHA-256: `a9e44459ac3d0c7ddbf02c158294e8e0afb0523abde25f7f95f4821b1b4c42e8`)
- go@go1.26.10/src/os/sys_windows.go (source SHA-256: `989bccb9cc276b2d3017d6725c28ae552767a78b8150e0e53139a84b1c4e2e46`)
- go@go1.26.10/src/os/types.go (source SHA-256: `ef72e4082700de60ba3e49cbe7669a11096eb387a5460497f1f8e5d33f7b0a8f`)
- go@go1.26.10/src/os/types_unix.go (source SHA-256: `d222d99ab7e873441136cd3daeb4c147f048ec5289338ed123f223187a90f285`)
- go@go1.26.10/src/os/types_windows.go (source SHA-256: `32bdbb028336c434525559642424295c71c0dd6e4fef7cec1fd0739ec80feffb`)
- go@go1.26.10/src/path/filepath/path.go (source SHA-256: `eb7d1ddb99b53810f008faf096d15cdcb28ae04c66024345fcfcd0b08397d93a`)
- go@go1.26.10/src/path/path.go (source SHA-256: `0a5024c364387eeea9ab5fe0b285f0d62355e8cd181cdbddfdfa227f6916d94c`)
- go@go1.26.10/src/reflect/deepequal.go (source SHA-256: `e4615c50434192fa65a53c6ccaaa060da744f67baad2778153015d9f4cb26cbe`)
- go@go1.26.10/src/reflect/type.go (source SHA-256: `d2c31904f1027eeca56fba323fcaca1e16d63f45e34c850e1f39c1790a4fb7a1`)
- go@go1.26.10/src/reflect/value.go (source SHA-256: `7b93df76cdaa42211f5568780cea85b0f4a117ce177def008370da20e6d7b237`)
- go@go1.26.10/src/regexp/regexp.go (source SHA-256: `addd07fbda0b622c063e0f5516b55c684d48ab9cc5d3a4a9f90d82b19c527772`)
- go@go1.26.10/src/runtime/atomic_pointer.go (source SHA-256: `364e60036abead10a020d763add8fe57fc493624c0b8f4a1c2e70de4edf21a3a`)
- go@go1.26.10/src/runtime/cgocall.go (source SHA-256: `9c4b5f7704610648df935176b35d942179a5c58e91ac22186ffa2786371ab8aa`)
- go@go1.26.10/src/runtime/debug.go (source SHA-256: `01d5f04213a154af89bee99c8acfafa16ff116b83765126ad970c55f4df1d222`)
- go@go1.26.10/src/runtime/extern.go (source SHA-256: `1afbe8b480e8c6b6a4d54b0f88f07606077a061f653198fa150e5d2a45e2bc8a`)
- go@go1.26.10/src/runtime/mbitmap.go (source SHA-256: `e107a955bb9db0f9c5d0d1f128978723f1e5b1c02c46a020f4310f47e8739bc2`)
- go@go1.26.10/src/runtime/mcache.go (source SHA-256: `9c27d457c7464565dccc3ad1900f0c177c57eec3cd35b7465c00971b6339da0b`)
- go@go1.26.10/src/runtime/mcentral.go (source SHA-256: `07fdadf8d90bc514805131fed9de5c8e00cd87580a6bf40f55f83353aae7247c`)
- go@go1.26.10/src/runtime/mfinal.go (source SHA-256: `49d40539afdd927355c9cfbfabc7130c561d0a19e68ca5e621082e58e9d96cd7`)
- go@go1.26.10/src/runtime/mfixalloc.go (source SHA-256: `4469bf858ac16f58c41adab050c5243d0dd8c5aad9d2e1c62b135635f7bc6f04`)
- go@go1.26.10/src/runtime/mgc.go (source SHA-256: `8ac2523177e5915edfaf4e122a903eceef08b5abd765f8d7ed96a8cc1ed3a504`)
- go@go1.26.10/src/runtime/mgcmark.go (source SHA-256: `68dac83a612229a11677b08cf7ee90101b3001096292cd80e890fb3a5cd4e817`)
- go@go1.26.10/src/runtime/mgcsweep.go (source SHA-256: `81139b39aeacd63e86bdadb6e921882e871bb30c87479b454aa5d82adc0c91aa`)
- go@go1.26.10/src/runtime/mgcwork.go (source SHA-256: `f475a9e19f9f361d76421b991bc331ee9f7c0b4d6ea375fa6a17c4d44bbd188b`)
- go@go1.26.10/src/runtime/mheap.go (source SHA-256: `6147548986e268b53fca158135149396a586cfd7d82aef7669b06f59a113bc04`)
- go@go1.26.10/src/runtime/mprof.go (source SHA-256: `cd26b885d64fa60668250491c3ea127d3364874059ed390d8cb9de570ecbb234`)
- go@go1.26.10/src/runtime/msize.go (source SHA-256: `0164d15295726608812842b429e241f553a314985a7c6aa48cb1ad24346a52dd`)
- go@go1.26.10/src/runtime/mstats.go (source SHA-256: `4a33458ef51ba74a46f63d66244c55d93356cdad9a2eafe7062c1f1800710ed4`)
- go@go1.26.10/src/runtime/os_darwin.go (source SHA-256: `0b87fce265c723dd50f4b1282422c17a5d0ce50bf115192e97e4f4a3844bc59d`)
- go@go1.26.10/src/runtime/os_linux.go (source SHA-256: `f211c96802c315ffa23343fc3a672f38713418f46658681da1f4d8187d9462a1`)
- go@go1.26.10/src/runtime/os_linux_generic.go (source SHA-256: `ed08e83bbb13787a044b1f12e952bc54f57085d9f5b4e80067c2c1e2d8527167`)
- go@go1.26.10/src/runtime/os_windows.go (source SHA-256: `d70b06f0c88d2cbc3d9534d30b1dc25e8a273443fdc4e6dfde6791a12ea4314b`)
- go@go1.26.10/src/runtime/print.go (source SHA-256: `28e2bbcd7e43215d955849acc00e9d65e7007ab12e061d853835cfb670ea2810`)
- go@go1.26.10/src/runtime/runtime.go (source SHA-256: `fe478d164989b5cd5783c9820d99879623780ba353f836f499238e76efb24fd8`)
- go@go1.26.10/src/runtime/runtime1.go (source SHA-256: `e9aea0b82b76931e2181563cc967112de12416a90c5467da8f2e178bbb04b8b2`)
- go@go1.26.10/src/runtime/runtime2.go (source SHA-256: `1aae37f643d6092ab3911ae5bf95b14941e88efc2c17609bb1db29f6c70841dd`)
- go@go1.26.10/src/runtime/select.go (source SHA-256: `93977178a5c507d2bd0bbd990613fa0c8a1d6d18b974fcbcd6c9edce6bd8b5bf`)
- go@go1.26.10/src/runtime/sema.go (source SHA-256: `451de314f2db754e318c0ab7857ec3b808d604845e6306eed064a57b8b2cd52b`)
- go@go1.26.10/src/runtime/signal_darwin.go (source SHA-256: `ec4119fe65a45c77833150bc11107e4a11bf1b062c04cd1c87c9731a01935d9b`)
- go@go1.26.10/src/runtime/sigqueue.go (source SHA-256: `4b645a4e6f3cdf78932195e8ff28f3b79b47bc4d6df7167acb19ade6c17366db`)
- go@go1.26.10/src/runtime/sigtab_linux_generic.go (source SHA-256: `0e618e00f7d446bc55e03083f76a276e481652603d10c96e7662e943d5101513`)
- go@go1.26.10/src/runtime/slice.go (source SHA-256: `ac595f1a822620fdc52271ef84f06b19c1f0abc0e53c0eb049dbad90b86abdbd`)
- go@go1.26.10/src/runtime/time.go (source SHA-256: `38c7dde4ad56b211cda6216c67e8f4c1b1f3a3ea544ff25b62fc83126e8a5f5c`)
- go@go1.26.10/src/runtime/traceback.go (source SHA-256: `b02b809e1bbb8af1fb1f7e81ef58962697be1be2591c5159d9a12959f59ead3c`)
- go@go1.26.10/src/runtime/type.go (source SHA-256: `abe23e7eab9ea0eac80b9b15ef89f0f7bcac887a56196b683ff6386026793ed8`)
- go@go1.26.10/src/runtime/write_err.go (source SHA-256: `fb40f33137552d750419154b3a19a346328eeac38dc9186f6450040ceaf89686`)
- go@go1.26.10/src/sort/sort.go (source SHA-256: `100e49103d23dd6b8335a50b17ae6e9db2b056cb1217d964dcf58852c3c6f8a1`)
- go@go1.26.10/src/strconv/bytealg.go (source SHA-256: `210ed5dd413d26575dabbb917c047ba52a144d133414f5e06b9595a4b30d3c60`)
- go@go1.26.10/src/strconv/quote.go (source SHA-256: `6602960ec8777de9926c1d02fbf5cb0828d4aefa38083865f4af17b6892c09a2`)
- go@go1.26.10/src/strings/reader.go (source SHA-256: `cceec92b1c65a8b3b2b5fc2991ea0a54b1bf6c5fb75981de5fa91c7b2dadbb6d`)
- go@go1.26.10/src/strings/strings.go (source SHA-256: `ba9478be775d11b6956c828ebf870a8f8e7b7318ee382834f1e741f02d1336ef`)
- go@go1.26.10/src/sync/mutex.go (source SHA-256: `3dec4264678744c4ce15a05ba0cc88e5e8b155f071b71d74eb7ed229ca72c114`)
- go@go1.26.10/src/sync/once.go (source SHA-256: `3bc00bab32acd134ddf50c87fdbd1dd70502f224b7684a801a077dfe0a05eb3f`)
- go@go1.26.10/src/sync/rwmutex.go (source SHA-256: `4ecc9b2fa7771cbe37e1670e4caefd7a28939683bfe1614fcf5e9d545261d17f`)
- go@go1.26.10/src/syscall/dirent.go (source SHA-256: `e46154ad02f16b9b50cf7aa7fe8625a39e4fb1b30213902ff6c07cf8f37e0dcf`)
- go@go1.26.10/src/syscall/exec_unix.go (source SHA-256: `026985e8e1e9be3bf272bad3f708ba600292e09275c0071a31b67f6c7907a312`)
- go@go1.26.10/src/syscall/exec_windows.go (source SHA-256: `031afbedd3b2e8ee58d890b8ff3df3df67f5556ff174768ec17e9fe31d44ae40`)
- go@go1.26.10/src/syscall/syscall.go (source SHA-256: `87f7daeec9dbf53bc2b2b2432f6014530218f1ca45ecc91bcb1c7e60858afd6e`)
- go@go1.26.10/src/syscall/syscall_bsd.go (source SHA-256: `b772fe50da97a5a4c76e480bf4240d3fdc85a6ee87513eb8f8a785812159d263`)
- go@go1.26.10/src/syscall/syscall_darwin_amd64.go (source SHA-256: `e8664a1456152679ccb827fab10e99bf19ee7c8c5359e4f8f5e910bcc3f50559`)
- go@go1.26.10/src/syscall/syscall_linux.go (source SHA-256: `e5fa95a3d18df64cf660e47d832c0fcf32fcc4681606aebcad6269ce522ef632`)
- go@go1.26.10/src/syscall/syscall_linux_amd64.go (source SHA-256: `e220a8ff4c78efdfca407b87d8db7dcb6b61e20238091dacfde1bdabccf7b42e`)
- go@go1.26.10/src/syscall/syscall_unix.go (source SHA-256: `31fd5823803feb47f4ea0f13d3e0db5bd3cc6e37467bbebd8abbba8467cb20a7`)
- go@go1.26.10/src/syscall/syscall_windows.go (source SHA-256: `f1b22dde6bb980cf036561cee7ece968daa260cb622f32c96cebd52477b7a346`)
- go@go1.26.10/src/text/tabwriter/tabwriter.go (source SHA-256: `202c00a1bc55b5f37a1ef66b4d9bd841fb713aa88801d881ee20f112e93591dc`)
- go@go1.26.10/src/time/sleep.go (source SHA-256: `9ee80a355252f5b71e9066175eea7a7c67e75724a483ab45258d6e5b5f023507`)
- go@go1.26.10/src/time/tick.go (source SHA-256: `c25a5f521f3925d21dedab88e8829cfed90b0fe6283e9c189d30a3f12942ed81`)
- go@go1.26.10/src/time/time.go (source SHA-256: `e8b6ea9e6e7cb890aeed1d24ee7163a2237b6fa78ffc564efb8688dba717f02a`)
- go@go1.26.10/src/time/zoneinfo_read.go (source SHA-256: `47e6cc315fbf5ba3eea95043fe7ec5e982d14b98574e60da22ac9b4f774f254b`)
- go@go1.26.10/src/time/zoneinfo_unix.go (source SHA-256: `d2001a2351a080f1a47ebb3999be74f99caceb4a8f1ea6369a3c843fa0442f08`)
- go@go1.26.10/src/time/zoneinfo_windows.go (source SHA-256: `831f1d99fc109507783841c65f0e196f063a1ce6490354869bf6df68aa626434`)
- go@go1.26.10/src/unicode/casetables.go (source SHA-256: `b70c5009086ffaef3d8191924481b9b2a39108afb42d9f2b979ecf670c77670c`)
- go@go1.26.10/src/unicode/digit.go (source SHA-256: `dfb14dbbcb918a6c7ee9be15293327d022f9d64a4c1a90f9c8049b682e799761`)
- go@go1.26.10/src/unicode/letter.go (source SHA-256: `b35a6d259313cef88d07e9e4027b8bd72c60ae140418f5be2c41a2ecea6074fe`)
- go@go1.26.10/src/unicode/utf8/utf8.go (source SHA-256: `3c145ee25fd7631ebbc670fd1236e8c163ca30ae8913c796489efbf3f267ad54`)
- go@go1.26.10/src/unsafe/unsafe.go (source SHA-256: `139f79859cb2c690b63b56afe5c385b694ca418ce2dead0ab9d6b1568fad933c`)
- go@go1.26.10/src/vendor/golang.org/x/net/dns/dnsmessage/message.go (source SHA-256: `e1cc3e8ed316b1b3e0153f5f3b6ee5ed2e0045d9ab3abf7eba7e00daadb773e6`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/runtime/asm_amd64.s (source SHA-256: `837d93ba923e8ee89144512a7a9c7a3aa54da91b0143fa149ac7732fe68756b8`)
- go@go1.26.10/src/runtime/cgo/asm_amd64.s (source SHA-256: `4592d4cfe237e01561d7a6329fdfe5c384a247afb0ff2da7a5ef2a355ac23db8`)
- go@go1.26.10/src/runtime/cgo/gcc_amd64.S (source SHA-256: `47427e195996d98726efea5dfb4093c8ef8887e96067095f0053f5adb62485a9`)
- go@go1.26.10/src/runtime/cgo/gcc_darwin_amd64.c (source SHA-256: `7e330fa22929e869933772a9f3c87d0b6880b9d680ca62ac8b7670532c2d6bbd`)
- go@go1.26.10/src/runtime/cgo/gcc_linux_amd64.c (source SHA-256: `d0f3b6f018c336e521368acc74c64d9aa8788f7587d434cb24d7ff9df73197ad`)
- go@go1.26.10/src/runtime/cgo/gcc_util.c (source SHA-256: `04d353f36b280597753a7c4a013d4b19b906227775f4ba7b2c54ec297ecfb504`)
- go@go1.26.10/src/runtime/cgo/gcc_windows_amd64.c (source SHA-256: `fd7a5ea96b01b68f320098a5210f34debb7025583194110bcd433d36333dfede`)
- go@go1.26.10/src/runtime/cgo/libcgo.h (source SHA-256: `68a44ce5a99a812c766944d8115f69d2a5f827486a44c00ca37451c160e93cb8`)
- go@go1.26.10/src/runtime/rt0_darwin_amd64.s (source SHA-256: `1c1215e5b8b8fa439f0239e5ac1eddb9f265c3a8950e67c889147ea6be7cf78e`)
- go@go1.26.10/src/runtime/rt0_linux_amd64.s (source SHA-256: `60586230585acae1beb5f9a821ac07381c3fd8b8c262f7fc8798fa3f7d0d854f`)
- go@go1.26.10/src/runtime/sys_darwin_amd64.s (source SHA-256: `11cb8d871cbd49f75550f9e4a89d349bb458ff243b9b3fba87bf964e7c337912`)
- go@go1.26.10/src/runtime/sys_linux_amd64.s (source SHA-256: `98ec6f0afedcc07d64f94c299c201b1ca5e7d52f7ed284bf6782c5ca7b9760f2`)
- go@go1.26.10/src/syscall/asm_darwin_amd64.s (source SHA-256: `f4ac337d025ceafeca0a412a4354afc695528ea5e492b967e34f0ba2bc7c64df`)
- go@go1.26.10/src/syscall/asm_linux_amd64.s (source SHA-256: `6c4b2e830d3918f72923c172e25c75c755b995d25b668fe26dc0a244311d8cf2`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/syscall/syscall_darwin.go (source SHA-256: `07eea080aae7b97b37943bb0964594a198c73b8ec35697b2b8912c0323cd129a`)

````text
// Copyright 2009,2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/math/mod.go (source SHA-256: `2ee912a367bbaeab29dbcf76f1112804b7eb543e8838e76bf869aab7eb5157d2`)

````text
// Copyright 2009-2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/archive/zip/reader.go (source SHA-256: `512482f33931f71f225fa76700847f31f9fc82cdc281338b82016818db2806e2`)
- go@go1.26.10/src/archive/zip/register.go (source SHA-256: `79c2b8c310f04ee33c86e8c7399fd73c3d2b6d81a7c6a711a6a740654ac3ff98`)
- go@go1.26.10/src/archive/zip/struct.go (source SHA-256: `0245d82aa77e2b9dde3ae1c75edab06f9153f801e0742dac0617f24c33d5609b`)
- go@go1.26.10/src/compress/gzip/gzip.go (source SHA-256: `9f5cef3dff42000c6df57de93429c10121c20afdda99b0c883fa60f94617a497`)
- go@go1.26.10/src/crypto/cipher/cfb.go (source SHA-256: `3d08078f78c09549271c047a26e112d75277c4713ef5366640e6a6360701e359`)
- go@go1.26.10/src/crypto/cipher/cipher.go (source SHA-256: `af744a7e24a704177bbd0bbd7cab39896583c5304aacb17bc7ee889adddc00a5`)
- go@go1.26.10/src/crypto/cipher/io.go (source SHA-256: `cd32813d82194c2c24144d0622cf3cfa95b2ce0dbbe39822a2b4ff7f3e38f5d7`)
- go@go1.26.10/src/crypto/des/const.go (source SHA-256: `fb1b9ff9539a17a3cffd374e75bfcedbebbaf23732365af56003f04697c6c3b9`)
- go@go1.26.10/src/crypto/elliptic/elliptic.go (source SHA-256: `5c265d23fd60baa8ea851a9463e981f9e34978696af0d8c60578b1fecd6762de`)
- go@go1.26.10/src/crypto/internal/sysrand/rand.go (source SHA-256: `76792623691f632274ea7996011bcfd90ddb52a83151f75fe149c255bb69cc1c`)
- go@go1.26.10/src/crypto/internal/sysrand/rand_windows.go (source SHA-256: `8122c2f83776b2e76f45631835d8038853886e5db728e24fe97c60179197ed8b`)
- go@go1.26.10/src/crypto/rand/rand.go (source SHA-256: `7c289648d453d0bdca0f05091d00bb86c94b69637080c748ff3d3476677521ac`)
- go@go1.26.10/src/crypto/tls/cipher_suites.go (source SHA-256: `6061822eaac77541fe38d46acdd4766eb70d2884f8bf84f3ce7f0471217ac73b`)
- go@go1.26.10/src/crypto/tls/conn.go (source SHA-256: `c5267d4cc54fef91f4346e82865bee81224599c080dd2c559b21e09dd1d7e7d6`)
- go@go1.26.10/src/crypto/tls/key_agreement.go (source SHA-256: `3874cd45e07283547ae3ab6ba4a10893323210ba153906fe586ce0158ba45862`)
- go@go1.26.10/src/encoding/json/decode.go (source SHA-256: `a91535f18f9a2c8380448b1c84c4c2017dbe3d530a10faf8babf03494042da48`)
- go@go1.26.10/src/encoding/json/encode.go (source SHA-256: `4012d78aa8236f3067a6330ddf2f719c7e59b262fcf4fff5e51b4abd1ccba235`)
- go@go1.26.10/src/encoding/json/indent.go (source SHA-256: `9ef9eac29fc4cc9d94bf305976891eb5e78a2644acd4f07ed7a0c186991ac908`)
- go@go1.26.10/src/encoding/json/scanner.go (source SHA-256: `2b16dd215274dfa8e0806b53cca144684d5b9c8a02319cf5ce10120c50e6eef2`)
- go@go1.26.10/src/encoding/json/stream.go (source SHA-256: `bb26b2442db6a686fc7f823ea9cc89babe26ac988ec7d8caf4a157aeb60d92f2`)
- go@go1.26.10/src/fmt/scan.go (source SHA-256: `ac416ad4fba7f9901c737a203950d5d57220d1401e4bc7f0a2420389b88a86cd`)
- go@go1.26.10/src/go/token/position.go (source SHA-256: `8ed202448629efff9959e4c2199330f4e9545707c63003813f55b1e33fdfcdbe`)
- go@go1.26.10/src/html/entity.go (source SHA-256: `4e6dd013ad9d909f4f8a93d200567c5fd88a2cc1d12b1c2046b3fd199b373ec1`)
- go@go1.26.10/src/html/escape.go (source SHA-256: `c3cf95cb15e3091bf25bc737711e14d252573812595d0d86e01192312df9938b`)
- go@go1.26.10/src/image/format.go (source SHA-256: `29605fb47ac733d1fe02476277615da520af04247933e38df606c4e513d4f1e4`)
- go@go1.26.10/src/image/geom.go (source SHA-256: `7cc4be83e978cc172e0ebb9713a6afd77ccf3ceb82d9cbc6f79082e698884ec2`)
- go@go1.26.10/src/image/names.go (source SHA-256: `7f0d1a350592934f7016ab1bfe6802af7d3fbec22fde175c936211ca08a60e3a`)
- go@go1.26.10/src/internal/filepathlite/path_unix.go (source SHA-256: `91b9b7d3719089cf1d6decc8cae667fe7ee32230b43d567def169d0ff66505af`)
- go@go1.26.10/src/internal/filepathlite/path_windows.go (source SHA-256: `ad313fae9fd1e150df082ad7f43c34b0830ce79433d0eef89c7724fdb07d06bb`)
- go@go1.26.10/src/io/ioutil/tempfile.go (source SHA-256: `50c19cc188598a4d483045b8d9ea7145ad2a1f241609dde5992fca5cee895f2f`)
- go@go1.26.10/src/io/multi.go (source SHA-256: `acb61a128d4e49dcc8b5dcd268ee1bd61008fd97af0e459f839166e3ae51384e`)
- go@go1.26.10/src/math/acosh.go (source SHA-256: `e1b20885bef3ccb7763ae868d23120194cf38faef70f1af45d2870ccb1a00140`)
- go@go1.26.10/src/math/asinh.go (source SHA-256: `7bd3753d8bb94eb86682c9215c6266efccd75a66701aad16fc94d9fbfb89ccb1`)
- go@go1.26.10/src/math/atanh.go (source SHA-256: `2afbcd791957596e7af481c3d23640f15e2d1149d8ed70fc108ddd1d5ab79430`)
- go@go1.26.10/src/math/big/arith_decl.go (source SHA-256: `7a556b5c1af9384f2fd1fdc4d548401e73e07a422df9fdcbf15ed4140b74a8bd`)
- go@go1.26.10/src/math/big/rat.go (source SHA-256: `c29ccef6bcf4b5f4cea3f4fa51aaf7fec169685b1bf831236f952532a6a901d0`)
- go@go1.26.10/src/math/cmplx/abs.go (source SHA-256: `b89dd35c4023b4bdae41d69955920e27fd36dc1da2534b8ad312428f6d94015d`)
- go@go1.26.10/src/math/cmplx/asin.go (source SHA-256: `10d60a9dc9aef2e0217e73bf648ae594c91e194de35dfb2eec96aff6e2a38b18`)
- go@go1.26.10/src/math/cmplx/conj.go (source SHA-256: `261caf2960f27907b1b48b5d7810bc79ff2ce1b2e5d2f07d89ceba51d6a7eaaa`)
- go@go1.26.10/src/math/cmplx/exp.go (source SHA-256: `1660975b05851b22747bffa9e1e8990c73f86a558e364c46b6b3054fdfc596a4`)
- go@go1.26.10/src/math/cmplx/isinf.go (source SHA-256: `30520c219f406b7900c62915396fb28fcf64f84640d99d855c1939477d0f4915`)
- go@go1.26.10/src/math/cmplx/isnan.go (source SHA-256: `4b636871557d31aaef791065b106132b1bb491b140f02c604c329593c2debf56`)
- go@go1.26.10/src/math/cmplx/log.go (source SHA-256: `4e78d54adda6eacf3ae1e5986e428b1b946d158657ef781b218694fb578cc85e`)
- go@go1.26.10/src/math/cmplx/phase.go (source SHA-256: `0bb918fa386b62b216d2bbcb3afb05fe56bccf5b2fb614e02f19c8be96406862`)
- go@go1.26.10/src/math/cmplx/polar.go (source SHA-256: `3e7db89588041d745b33866cf8938e4d061b1d6a50fae682d64354281c13a81b`)
- go@go1.26.10/src/math/cmplx/pow.go (source SHA-256: `a64cf1d1955098bfec0b2da0cf0122ca5777e33c9f623fcc03f49d950e90e3ca`)
- go@go1.26.10/src/math/cmplx/rect.go (source SHA-256: `a4d26bfef1c800ca9176f7b349e72524b9602305772e667a6749abba6629802e`)
- go@go1.26.10/src/math/cmplx/sin.go (source SHA-256: `de24e3dd23fea2fd71d5206136f4c7786568d785790af1cfae01ef0002c821bd`)
- go@go1.26.10/src/math/cmplx/sqrt.go (source SHA-256: `d55bfa6adb2448163ca500b9f36175e18915c6d9973a7fb117e6c6acc76aecc5`)
- go@go1.26.10/src/math/cmplx/tan.go (source SHA-256: `0afd68f9317f3a278936d7ac50fe0614a6bee7d1a7a0023ffcc7caa25ba1c660`)
- go@go1.26.10/src/math/copysign.go (source SHA-256: `9d7e6179f8bf2df156b1bd054de7ef12afe26205c1b786e3ec059b1a80284bda`)
- go@go1.26.10/src/math/dim.go (source SHA-256: `d9861a32bf027b42a280c161ef764794c8437ec028706e8db316bcf56cd023e0`)
- go@go1.26.10/src/math/erf.go (source SHA-256: `31cba8a8f66d6081eb393b85e6a0d2a3737ebb3fc3fb192867e393a2907c0c8c`)
- go@go1.26.10/src/math/expm1.go (source SHA-256: `b05c184f92e2e4bdfb09eb34e6d2474fe4570a4eaf5eb1e128d2b6904908d63c`)
- go@go1.26.10/src/math/gamma.go (source SHA-256: `627e52bb7a984ab2da61da481729f317eae7a320a2980c72aa6077ca9e6bd03b`)
- go@go1.26.10/src/math/hypot.go (source SHA-256: `0118c2cfefb7f3085b6267439fff271843a06eed5445289c260b504c4954d6ed`)
- go@go1.26.10/src/math/j0.go (source SHA-256: `8d226442c38eeb79e5273980fc5d94f0ba9d4c4407d0e9a7e8db45838fc128f8`)
- go@go1.26.10/src/math/j1.go (source SHA-256: `599ddf54522dd0cf19e3cdcee1298a5ca98eaa529e6e647fb8311373f5bef1fa`)
- go@go1.26.10/src/math/jn.go (source SHA-256: `b588c5389aa3ac04eeb72b29a610bef4af02b8ee7faf3d50717d2da2219bc144`)
- go@go1.26.10/src/math/lgamma.go (source SHA-256: `c07c6b7b6d57e97f47531abca86accbd23e0a0d49d8d5cb60c226e87df9c6804`)
- go@go1.26.10/src/math/log1p.go (source SHA-256: `10cb153951ada1c8c5ed3897caded00d90a7314c0640812847c3064a0c5f4b06`)
- go@go1.26.10/src/math/logb.go (source SHA-256: `9cd63343d6d38130520af9e1a107b5bce1e670357a6292e157fcab90c7b82365`)
- go@go1.26.10/src/math/nextafter.go (source SHA-256: `f52cd1442fadfca48ff5f3ea8620293655180486e3ad1f88288568c797f29e15`)
- go@go1.26.10/src/math/remainder.go (source SHA-256: `b86740c9a71e4144feec5e728c13f18b157744151d55eeaada90c9c9f26cefef`)
- go@go1.26.10/src/math/signbit.go (source SHA-256: `4680ea79d079fdf40773678c72027558ba6af61d649fe51e73d1214bb905783f`)
- go@go1.26.10/src/math/sincos.go (source SHA-256: `d1a9614979455110d808aa194fa6e2451cbfe837b11584311fbb356a664a52e5`)
- go@go1.26.10/src/mime/grammar.go (source SHA-256: `85d1052199201e6ad52e27db50104842295dbf10e87f3914c803c6ce533772a6`)
- go@go1.26.10/src/mime/mediatype.go (source SHA-256: `2179fd195097a382a2ac3cdcbcbc6961f2770dbd61adfabfc60535c1acc5faa0`)
- go@go1.26.10/src/mime/type.go (source SHA-256: `6ff51e45cc2357108cdefb4665586755d67fb1e6467907caa2afeb94711afeb4`)
- go@go1.26.10/src/mime/type_unix.go (source SHA-256: `ecf9e4a99edf73edd2e9088ca0ba0c8ec99a094188474536e59dc0a921a850ca`)
- go@go1.26.10/src/mime/type_windows.go (source SHA-256: `f893579a237fd8373f4b4d576af294bb12f62fdc8aac952ffadbae0906ed6cf6`)
- go@go1.26.10/src/net/dial.go (source SHA-256: `9963e94645bb8620e61bf219236a5d859b17de161344d433b93bedd6f3ad1e1e`)
- go@go1.26.10/src/net/fd_windows.go (source SHA-256: `8008f4a0bd9fc58a23b3fea17ffdb15e4739c19585c60589b764076d225972ff`)
- go@go1.26.10/src/net/http/header.go (source SHA-256: `b4fa959c3db03d4ca67a487b3cbb25dcefecea45ca3dd10fef13cfc32fa6d80d`)
- go@go1.26.10/src/net/iprawsock.go (source SHA-256: `fcd9c09b13588715fbfe05b9dcf381623d3b3241dfee39c93aa2044421ac16e6`)
- go@go1.26.10/src/net/iprawsock_posix.go (source SHA-256: `f2ef8c16772021695b2380ba5a3dde6e4a16c953ec0236b7a92d3619483c811c`)
- go@go1.26.10/src/net/pipe.go (source SHA-256: `55d3e00ea0c68e17e90da06eaa26a97908776f9e55778f12768c8e0d9c6aced1`)
- go@go1.26.10/src/net/textproto/header.go (source SHA-256: `8ff6b9ab3af4f79aead3b4a3e07aa11e124980b0ab8bf24859fcc7acd2cda7d3`)
- go@go1.26.10/src/net/textproto/pipeline.go (source SHA-256: `d8b78ffe222aeae796fd97b6feb8eec481fb5325bfcda1448e1f6b1866557e83`)
- go@go1.26.10/src/net/textproto/reader.go (source SHA-256: `fe649123b4df42a4aeb7713aa3f8028de61897072874f37c67c497bd6bbf734f`)
- go@go1.26.10/src/net/textproto/textproto.go (source SHA-256: `e4f0d94752345b5dec4fe7f4805531e8036a9cd44131d807cafa700dd5323901`)
- go@go1.26.10/src/net/textproto/writer.go (source SHA-256: `88ff154ebc598835a777832e72e013a6c2280f70fc74a9af515b80cb6ffd0907`)
- go@go1.26.10/src/os/env.go (source SHA-256: `e7eb15de93fccf8c0706285c350f32eab6c4a1f752bb8594ce7da3a0215ee92b`)
- go@go1.26.10/src/os/exec/lp_unix.go (source SHA-256: `5ee781400e1f04ffea869496f161555f31d1af0ca521fb77ef16f54522808746`)
- go@go1.26.10/src/os/exec/lp_windows.go (source SHA-256: `b43699d7c9c5fa2cfe82edb55be04b509619232ad514a34a11c481d43290a6fd`)
- go@go1.26.10/src/os/tempfile.go (source SHA-256: `a337c2a6dc2e4368b2dc0015ea7c0b935b68b14885a0771a6cbe27e588c2cd16`)
- go@go1.26.10/src/path/filepath/match.go (source SHA-256: `ed53eafdde82d19db7c2747c514ec9f7d0b51458d550e707167a907376050e96`)
- go@go1.26.10/src/path/filepath/path_unix.go (source SHA-256: `d92b04d88e202d2916f803ceb3c8806475b510907a142571f4e9ba70b52c1c1d`)
- go@go1.26.10/src/path/match.go (source SHA-256: `0c6f2ce8e03bd945c71e669e129e65c555917192e054b962e0ca7c473defe959`)
- go@go1.26.10/src/runtime/cgo/cgo.go (source SHA-256: `619406c7dc9470188f820c53ff02e6372459e533b7cb75c717dd2ba6b17dd4fd`)
- go@go1.26.10/src/runtime/cgo/iscgo.go (source SHA-256: `eb9cb0774da854424967b8b1883efa657c166f2911091722356a9f986f1ee8b1`)
- go@go1.26.10/src/runtime/complex.go (source SHA-256: `acb03107a1de5df055c14f0fefe0a5c8e33ff6658445d10e1eb23e8e99e1d8ad`)
- go@go1.26.10/src/runtime/error.go (source SHA-256: `138957e6d2a126af0d6b2d032466fbe2bc3c95bb57b058a00987e814b2385e1f`)
- go@go1.26.10/src/runtime/mem_linux.go (source SHA-256: `08cc392d6b540338b966c9593691a700d8ac99ce4081f444f792b994e37b52c7`)
- go@go1.26.10/src/runtime/mem_windows.go (source SHA-256: `1f915e8b6034c4815fcd52ba9e878808a66f81dd196a29d4acf8e202ac124837`)
- go@go1.26.10/src/runtime/pprof/pprof.go (source SHA-256: `d8c92bf67507fa0385b50d848dcf4d278922e9acb1755ae7d6968baf82466cfa`)
- go@go1.26.10/src/runtime/softfloat64.go (source SHA-256: `94f52048c124deff27c7b4c0215eab39700b16cd9a3781e6453f7900a4eb2005`)
- go@go1.26.10/src/sort/search.go (source SHA-256: `78418580ae0329bd5e5b33e91223f62c42010a3d16c17f37ff953ac094611686`)
- go@go1.26.10/src/syscall/env_unix.go (source SHA-256: `7a75aaf2e881a7d48cdb449c5854925dd81dabe45cee6c54e508727b85c7dab2`)
- go@go1.26.10/src/syscall/env_windows.go (source SHA-256: `5fc24ad0663552b191c8f1014a5cecb976855778b1d7e57244040fbb74c0df32`)
- go@go1.26.10/src/time/format.go (source SHA-256: `f288c56858ab25828551681778f006b591ed140fc12e60122b6a9a7e35b07eb7`)
- go@go1.26.10/src/unicode/utf16/utf16.go (source SHA-256: `6b8b3037015f41351d4a323513b57b8ac90caa69f69be31eaf02ca2a4632d0cd`)

````text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/math/dim_amd64.s (source SHA-256: `5569939ae83d71a6edd1503908f46b4480d7cca4069a59a9eb620ded1d724e90`)
- go@go1.26.10/src/math/exp_amd64.s (source SHA-256: `6b0b9ca5b73d27e457fa3ce283dec9fdbc8411ac3b28afe262059330c32e9980`)
- go@go1.26.10/src/math/hypot_amd64.s (source SHA-256: `47046a6828e174a9b0d6fa4607d64602215123f2304b1592fd28635722ae5bf6`)
- go@go1.26.10/src/math/log_amd64.s (source SHA-256: `5764160fdc45d72cc258a67034fefda930857a34d68d667dbb0bf73fd0e0e88b`)

````text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/mime/multipart/multipart.go (source SHA-256: `edda61c98001eecf78e48f87c709f0a2c9a3972951a04b48d15724d9ee8df9e7`)

````text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
````

### Original source notice

- go@go1.26.10/src/archive/zip/writer.go (source SHA-256: `ea7d65bc45ab1eb96c72fbbddd31c09c510eed9b9d2ac33b535b7915efed08ca`)
- go@go1.26.10/src/crypto/cipher/ofb.go (source SHA-256: `3c91daaffeb017c010183d6a08ec7069f68394c6a28d328c08eaf4c097cfcea4`)
- go@go1.26.10/src/crypto/crypto.go (source SHA-256: `c25f42ceaf88072a813f73ec185db3a4e2ee1b80293abe28edc119e09c73df33`)
- go@go1.26.10/src/crypto/des/block.go (source SHA-256: `b2db89bff122de6f8ad47510522c297c5e7adbbb8527543a9beb50c303985c92`)
- go@go1.26.10/src/crypto/des/cipher.go (source SHA-256: `0dfa5c8ef6e6043f5031b412cb66f58b8ce197a7a3494c0effc475a7663a0335`)
- go@go1.26.10/src/crypto/dsa/dsa.go (source SHA-256: `365b396f4541b239b03ddb8370eba2e8209c51ad4463530cecfbee5b30ea69bb`)
- go@go1.26.10/src/crypto/ecdsa/ecdsa.go (source SHA-256: `61d06e6707ef461f7e77fcff654430dea46970edd446ac49dc3692dd5b43ac55`)
- go@go1.26.10/src/crypto/rand/util.go (source SHA-256: `109920499d49126f14b7d78469d4863035da5c5ad3e4737bf5056a9123fb172c`)
- go@go1.26.10/src/crypto/x509/cert_pool.go (source SHA-256: `d995d6e88af70f36a345185420bac88c58e86ebed1b3eea8e087228eaa7da03b`)
- go@go1.26.10/src/crypto/x509/pkcs1.go (source SHA-256: `d494f7d49f03daa26fde42128cf64df1ce737fe58625271f3bc195791b6b2d67`)
- go@go1.26.10/src/crypto/x509/pkcs8.go (source SHA-256: `6b62c951d4fc9ce980f3f47f58c5408c265d56e82f5ded3fa9e61911baf64b20`)
- go@go1.26.10/src/crypto/x509/pkix/pkix.go (source SHA-256: `7b8cb67e7e0fa0f067027a09a6d7793507d9ed74c9c20265c7514dfc63998d9b`)
- go@go1.26.10/src/crypto/x509/root_unix.go (source SHA-256: `421c066f193250dc1adaf13fce4779fc3dc576e20faa86e43bc22af6c9536a14`)
- go@go1.26.10/src/crypto/x509/verify.go (source SHA-256: `238eb6c7d7fc7ff6462cfdb0a7d5da4e556de4c7443a66e0ff79dbbe4c179081`)
- go@go1.26.10/src/database/sql/driver/driver.go (source SHA-256: `6510995fa991f38da3728bc7154e89f358abfca8dfd1be3cb1bc1c127e57d4af`)
- go@go1.26.10/src/database/sql/driver/types.go (source SHA-256: `e7dbe9c23d3b26ed0d167bdd9d6c46b0eb3c053c127baa3155793fd83b2c3ef4`)
- go@go1.26.10/src/encoding/binary/varint.go (source SHA-256: `8dfcb2f0eb8b854d275a0f8facff7ce68b206a56e7d3a3bde0a379c2433f5726`)
- go@go1.26.10/src/encoding/csv/reader.go (source SHA-256: `c37960d74220590547febce08339aa321f4d85c524740451e14692ebbc9fef1e`)
- go@go1.26.10/src/encoding/csv/writer.go (source SHA-256: `a0a4bdd01fddd75efc583f822599118320d30b8269f05eb957777c28e2aec1ee`)
- go@go1.26.10/src/encoding/json/tags.go (source SHA-256: `e7300718fc8862f4dbc8bbdde33e4a0335124bb028e68062b52d3c042c374d6e`)
- go@go1.26.10/src/encoding/xml/marshal.go (source SHA-256: `b91fb2437083333e4b35cc43f47377202e177de37fc7b1b0e3bd656dc400fae8`)
- go@go1.26.10/src/encoding/xml/typeinfo.go (source SHA-256: `3c71ff68c7c60f34cea939e6e1c851b2db82864a691f447437ef6b082496339f`)
- go@go1.26.10/src/errors/errors.go (source SHA-256: `865b9cce008bf68d4d2cc4146455a9c5692756bbe57be29ee1ef0de345dafe1e`)
- go@go1.26.10/src/go/token/serialize.go (source SHA-256: `79743e48958a1e72a49a56b09b3546db6205ee858bf9af913a4ee0665657e16e`)
- go@go1.26.10/src/hash/crc32/crc32_amd64.go (source SHA-256: `548dc603730f2b153b1568bacb155ac57568e8b474d5be371971e4f785ab2865`)
- go@go1.26.10/src/hash/crc32/crc32_generic.go (source SHA-256: `ab955ded37b50660ad89f448fa2d6129e73cee97e4755770c0c695d666fb065d`)
- go@go1.26.10/src/hash/fnv/fnv.go (source SHA-256: `c954d2dd054c303c3221095fdeb863228ac4f4824efca330fd4babb99a074b91`)
- go@go1.26.10/src/html/template/attr.go (source SHA-256: `29ec0f227a963c1a3244d3c8fc8508c72be2c89c882d8cb6b6acba01b2647a12`)
- go@go1.26.10/src/html/template/content.go (source SHA-256: `cc62ea422ff9ccec1cef0fb3bff7eedb8ea3a0ccfa72706f75cbff97394086e3`)
- go@go1.26.10/src/html/template/context.go (source SHA-256: `0a5bb059d9b7f780904c987b8dcfda3ee7e781b7c6b269862f28ec0eb4931241`)
- go@go1.26.10/src/html/template/css.go (source SHA-256: `52d3584d9de08ddd1c6cde00a5b08fda956f96f4e28ef45be45737ae0869e95e`)
- go@go1.26.10/src/html/template/doc.go (source SHA-256: `235319c200e7fd0ecd4da6edd435a68e41e31a8c58ab5690387c9562a5a2e11a`)
- go@go1.26.10/src/html/template/error.go (source SHA-256: `7c0e67172971c75b7620f08a6609552ab12d5b7e6418e9d1da904f1843541b36`)
- go@go1.26.10/src/html/template/escape.go (source SHA-256: `003ea3bf183b8adf0a942cec4983549a8d14d3f322b95e607b06e8b89e4579aa`)
- go@go1.26.10/src/html/template/html.go (source SHA-256: `7092c5550650f25fe6fd668036839fcf5f2d46cf7ca77f7903ac7f8086f40137`)
- go@go1.26.10/src/html/template/js.go (source SHA-256: `869c440188ce7d12560f112aae6ae635888949ce7ea92e97bb63a0a11936d47a`)
- go@go1.26.10/src/html/template/template.go (source SHA-256: `a1289e3b6a81c9cf662459de2865c54cc8302ac1e829bb3a5949fda33c7016a4`)
- go@go1.26.10/src/html/template/transition.go (source SHA-256: `d0ba111b4ba9bfdf7c929a3abe17bccf10c0f97c207c3fa6c30bc507e8169887`)
- go@go1.26.10/src/html/template/url.go (source SHA-256: `f23df5e614504d215402412f7059750b80a77ff4b97a298db1e30377a78f3a17`)
- go@go1.26.10/src/image/color/color.go (source SHA-256: `03511ed250d0b4dd20605d5dc3465164e559033a6ff7e5cd8ac89d2e19c3ab6c`)
- go@go1.26.10/src/image/color/ycbcr.go (source SHA-256: `7afc91f21cb8eb68027f54180a862137bdef26cfd59389ecb609c0f87923de0c`)
- go@go1.26.10/src/image/jpeg/writer.go (source SHA-256: `d05aebe5e957345bc743652a7fc4570b183c688500c0ae871bd171dd3ad40082`)
- go@go1.26.10/src/image/ycbcr.go (source SHA-256: `84f5be8b40f9b5442d8c00596ba42baaa875cac1c6ecefb8c1217f90b100ecbf`)
- go@go1.26.10/src/internal/poll/sendfile_windows.go (source SHA-256: `9a1415e670597dd79034a62b4038985ca6ce63411c8978b5861ee91943df0beb`)
- go@go1.26.10/src/internal/poll/sockopt_linux.go (source SHA-256: `88fe3b2092de02061c80abb772023eb6330922c37205fd0631c43f082d9d4151`)
- go@go1.26.10/src/internal/poll/sockoptip.go (source SHA-256: `5c998151901c935bfc82619592cd5e2d559869c1dddf8e15822b64bc67e94c72`)
- go@go1.26.10/src/math/sin.go (source SHA-256: `913640c10f5093e8cdaafa1ebebbf5175b5a156bdc0fa89d77350dd4dad3be9b`)
- go@go1.26.10/src/math/tan.go (source SHA-256: `b1a549f6d267acf37399b29827ecc4cc152cfb655d8a2ed6dfe56a67e58f8e60`)
- go@go1.26.10/src/mime/multipart/formdata.go (source SHA-256: `cf2fc05db29dc4556ae12719c4db59b2f90b681eef5e74c5229a3c00e3a5bd92`)
- go@go1.26.10/src/mime/multipart/writer.go (source SHA-256: `efff79ea10e83051650936eb649416aeacb97baa94d62322f264e4cfcbe47c39`)
- go@go1.26.10/src/net/cgo_linux.go (source SHA-256: `f4fd072f69752ab9cd08f3117edf2b0b18bb1860030b78c02e740c6f0b1da643`)
- go@go1.26.10/src/net/cgo_unix.go (source SHA-256: `54e949492547503416f6a2f411f102a13b087769b2dd7b661fdf88a7a3d6531d`)
- go@go1.26.10/src/net/file_unix.go (source SHA-256: `e7a51efd33c3fa150c901dbac835c686321f331a3d73ad9f02efd9340f0f828e`)
- go@go1.26.10/src/net/file_windows.go (source SHA-256: `33005305d01549b938d1f6c08fa7b3ac0dd899f4575ebf9cf9522ec1cbac23e8`)
- go@go1.26.10/src/net/http/doc.go (source SHA-256: `5c88441de20755d855ddd90b47d7b400a0b46b5bcd18e731f27c03cd89874634`)
- go@go1.26.10/src/net/http/filetransport.go (source SHA-256: `12d958d23c97cbea79349df390626c3b8c1afdad0c2f626a1f5ec584959ad00b`)
- go@go1.26.10/src/net/http/httputil/reverseproxy.go (source SHA-256: `ba027076fd175fb0cc98dd14abbe4a3977882faacaaaa645c4dd25460067708a`)
- go@go1.26.10/src/net/http/jar.go (source SHA-256: `ce89dafcd7d38806241bc002787c159a578f50097aa23e9e88e3055997356b0c`)
- go@go1.26.10/src/net/http/sniff.go (source SHA-256: `3cf7e02214e5c8eeca9f5b12940c17723f4cd7a243c33b98386849c918644667`)
- go@go1.26.10/src/net/http/transport.go (source SHA-256: `643bce4b4fb016ed57c18e9d087e6c298faf3263d83728107c6b5cdaa398fe87`)
- go@go1.26.10/src/net/interface.go (source SHA-256: `c9200d071607df43e3759234c213a1385f24fd5582a03dde85ffecf7b649af0e`)
- go@go1.26.10/src/net/interface_bsd.go (source SHA-256: `bf0b163a34df03fdeca2f3f3d502d5db0c618ce7de80f22c8705871cefda5786`)
- go@go1.26.10/src/net/interface_darwin.go (source SHA-256: `ff7b69721be8d54b62eb80e894e32c547b1d8d2c4488f97e4e7b49949a428adf`)
- go@go1.26.10/src/net/interface_linux.go (source SHA-256: `e56dc3a6a195eeb374c9269b5684332f2cb0655cc37e5f5a900ed6e5d7e8d085`)
- go@go1.26.10/src/net/interface_windows.go (source SHA-256: `e02ba97662a41f6d3a163f8c565537194e72b447b12c345f76bc82bcde5e42e4`)
- go@go1.26.10/src/net/lookup_unix.go (source SHA-256: `7fc0ecb91aa268d3dbd6dad3b6b806d92223e11782d6fc1b4dcc4f5f9a4b788e`)
- go@go1.26.10/src/net/mac.go (source SHA-256: `01e3450adfec9610ec82d1310234a8ff640c0d3b606fb9e4aa8b4c9c77cf056a`)
- go@go1.26.10/src/net/sendfile.go (source SHA-256: `ec2efcdccf41435ea32a492c99c7f49c650607331d9d608ff39b543d9f4870ae`)
- go@go1.26.10/src/net/sockopt_bsd.go (source SHA-256: `1915e013449314181275eef5ed0316dc9a7a78209884e500c5128dc04cf7636f`)
- go@go1.26.10/src/net/sockopt_linux.go (source SHA-256: `06d62bcd95cf4f26c77f1c1dd311dd4785b042bbd9c0bc77e8cba6be50df70fc`)
- go@go1.26.10/src/net/sockopt_windows.go (source SHA-256: `fa834673d71b544925baad82ef36d686486e892b857ba6c4d7bd915fbcdafd15`)
- go@go1.26.10/src/net/sockoptip4_bsdvar.go (source SHA-256: `8612f7106dfc4bb73b61d5514a3854fcb99f387407feca59c905320377209dd4`)
- go@go1.26.10/src/net/sockoptip4_linux.go (source SHA-256: `cba1d46fa93c4a71ab11077e03eca7a85f2279ca7456bc38577cda0115cc7948`)
- go@go1.26.10/src/net/sockoptip4_posix_nonlinux.go (source SHA-256: `96571c040eae4a44663d097405eed6dda384e266ca8d778597f4fea2946493fc`)
- go@go1.26.10/src/net/sockoptip4_windows.go (source SHA-256: `ce517ee35f335e3025019a47ced7b4fc51f55ee9ea7bfb88a35e34113cc610f5`)
- go@go1.26.10/src/net/sockoptip6_posix.go (source SHA-256: `c43c1559e5ffafebb1800cd8f88db067fd2d9af1c9759176d3573a8010b31628`)
- go@go1.26.10/src/os/path_unix.go (source SHA-256: `8a5d4913c2ffd085b54df516e983916ac56c660de9123c0877d57788267861c1`)
- go@go1.26.10/src/os/path_windows.go (source SHA-256: `7256e7698c7aed7eba0c3e5e3dd441c7886dfbe21864734f3f8bd3a549cbff58`)
- go@go1.26.10/src/os/user/cgo_lookup_cgo.go (source SHA-256: `9ed2800b5271dcaa260b4d335dc0b65c45c3ac7a3dc359ce69126f660a19035f`)
- go@go1.26.10/src/os/user/cgo_lookup_syscall.go (source SHA-256: `9a6ebe71d0237261356f315105940a0c3bfe58ce1e802870263cbbca9f5a4fdb`)
- go@go1.26.10/src/os/user/cgo_lookup_unix.go (source SHA-256: `858bb7b3f4d6c2c1adb29a8528dd025f7f151a61a039bfc5fb1933e5e3c35075`)
- go@go1.26.10/src/os/user/lookup.go (source SHA-256: `e93e23d254d43872f9d14fcda58ddd8999791b8b22471180c79a46c8238228d8`)
- go@go1.26.10/src/os/user/user.go (source SHA-256: `e8e9a60b9455ffc3658b861789b525b4d1e85ff91d0af992eeac662b10ff91db`)
- go@go1.26.10/src/regexp/exec.go (source SHA-256: `cc9fda63ad7056d4c4abaaa49c00470f726be35ccee53f1c3fe73259e9634dd5`)
- go@go1.26.10/src/regexp/syntax/compile.go (source SHA-256: `2e1f753be77bf1fa54aed2a7f98932394324b82c91109e19d93283f482f7ee01`)
- go@go1.26.10/src/regexp/syntax/parse.go (source SHA-256: `d0f657113901da39d6118e4172206a0358b71c96ef43051cc0a741a0a4d04820`)
- go@go1.26.10/src/regexp/syntax/prog.go (source SHA-256: `7a24ebd22eec504ceaeffbd33641ae330d3ca666d54f6a29ce6d96fc127e6c64`)
- go@go1.26.10/src/regexp/syntax/regexp.go (source SHA-256: `72fad4645d96cf8b0ef0ea4a9dae21f09e7bea090d1bdd500732d811982b2192`)
- go@go1.26.10/src/regexp/syntax/simplify.go (source SHA-256: `d22b48ecde668bc22cd8998492c752ec088da5d027d61b7f04f7a9298796eec0`)
- go@go1.26.10/src/runtime/cgo/callbacks.go (source SHA-256: `133c0935662691a05b8cc9147d6b91ce3ec780d13942a06c322e122836ba6441`)
- go@go1.26.10/src/runtime/cgo/setenv.go (source SHA-256: `74970e37595cf53adc4971d53e68c713fd198a37cc170c7bcd75cb10ce63a1a3`)
- go@go1.26.10/src/runtime/cgocallback.go (source SHA-256: `82d38f6a32c4f8763af62ee77eb3aa6e8f52ba9c2077e7491b109e6530c2b067`)
- go@go1.26.10/src/runtime/cpuprof.go (source SHA-256: `d90aebaa1a719f339fde318cc8ec98e7312824ec2e2ad469a84fcbc6b5511c5e`)
- go@go1.26.10/src/runtime/lock_futex.go (source SHA-256: `c411429c46c093a5214ceebb3557f66de31056826fc2ab234ebfae70c0a2efb5`)
- go@go1.26.10/src/runtime/lock_sema.go (source SHA-256: `75d2a07f05071f1537a1fad231e3397ae33ebac5026c66d09959b7c898f1ea5c`)
- go@go1.26.10/src/runtime/signal_windows.go (source SHA-256: `19742597c5c6375893852292cc846127eac45759e9e397e8c10bb897fb470397`)
- go@go1.26.10/src/strings/replace.go (source SHA-256: `fbe8467c532ea7e5cdb848577423c4f218cbe2776907fb5c8c091006c786cdba`)
- go@go1.26.10/src/sync/atomic/doc.go (source SHA-256: `fb448dacfc6b4d5c74ae203224dc31088c2cf45fc42391a4fdb04a47ccc62413`)
- go@go1.26.10/src/sync/cond.go (source SHA-256: `324314fd3a484897ce540a242f9de85d00407040c8cd2ee7db5d5761f5d14d4f`)
- go@go1.26.10/src/sync/waitgroup.go (source SHA-256: `274e43a7cbc76eb0da4635a2680785200dfaf94d85e2d55cdd5122c76a4791d9`)
- go@go1.26.10/src/syscall/dll_windows.go (source SHA-256: `fea84c96a3b8cbe83e7670fdf87db829de71e192d099d9e8ef1539dd9991a5ca`)
- go@go1.26.10/src/syscall/exec_libc2.go (source SHA-256: `f44f5aedfe6355026d8e5deaf8b2a48eb7d6f17f33f009f428e04bbdd9057a13`)
- go@go1.26.10/src/syscall/exec_linux.go (source SHA-256: `566314c12ef25bf80312c654b58648ce6f94cb427b8097b4b2f10993058d4366`)
- go@go1.26.10/src/syscall/forkpipe.go (source SHA-256: `f8ffc4fdc8ae2bba416435ad4a0b89fee12aba4d65fcdcb2ace43b3a8da5ea96`)
- go@go1.26.10/src/syscall/lsf_linux.go (source SHA-256: `a0cc8b45b9996b8fd1bc334dd81051dde04357edc4a1c8ea68e5b539390d927a`)
- go@go1.26.10/src/syscall/netlink_linux.go (source SHA-256: `569949cb048d294b11e02078525d4ac2be2b6d8520575fe4c9c3234ef418242b`)
- go@go1.26.10/src/syscall/route_bsd.go (source SHA-256: `be9fb9750f84e730d9b563bb5c93d777191c446d43e5c6283a8924c3eb5f707d`)
- go@go1.26.10/src/syscall/route_darwin.go (source SHA-256: `f2eaa25bfff546087a3c7d9667635598c149e05991a59d82e8b15770fe056b8c`)
- go@go1.26.10/src/syscall/sockcmsg_linux.go (source SHA-256: `5a8ebab526ca05814dd14aa8970b6f177e484cd666ce36ef6be8eb4c502fa93d`)
- go@go1.26.10/src/syscall/sockcmsg_unix.go (source SHA-256: `8e78ac0374da70229837cd41eed4fd3a5ad06dbc90ef3ea4136ccedacf762616`)
- go@go1.26.10/src/syscall/types_windows.go (source SHA-256: `8f1702a7af541a173bc3687a9f82a887daab49b4ec61d3eb4ca9f4a287fff2fe`)
- go@go1.26.10/src/syscall/types_windows_amd64.go (source SHA-256: `0ba5eacd9b4f43c6fa7c5789fe132da98668fd07c9fcf217e1435b8561376c59`)
- go@go1.26.10/src/text/template/doc.go (source SHA-256: `0f74f7ac85529540095185d9905584606b9140a3b0a1ae7b3e064df3470387d0`)
- go@go1.26.10/src/text/template/exec.go (source SHA-256: `ca14feb551ec536deeddfd46efdfa60d1c328e6d1b00cd5aa6ca944e379c282a`)
- go@go1.26.10/src/text/template/funcs.go (source SHA-256: `6cc9d22f1a20a993f38810128bcb21f0bdd2d1bf612f4a279107773b7e207f30`)
- go@go1.26.10/src/text/template/helper.go (source SHA-256: `a9694b886bed5ed1aab68ff030176b01e3cff6ca5a5c36d8b14c77751ca71d06`)
- go@go1.26.10/src/text/template/parse/lex.go (source SHA-256: `c9af8ace95f3bb1d8cf13e1bfbbeb72cf9b621acf09f7c000a2b9bd468ec2dd1`)
- go@go1.26.10/src/text/template/parse/node.go (source SHA-256: `22cfac27972764f663ad210af96ebd4e69c51e40088f074b5bbc9db6edefccb8`)
- go@go1.26.10/src/text/template/parse/parse.go (source SHA-256: `b1e442c36924555be5f3b31e88d683991b1e94d4926d20269a7536e5be2df584`)
- go@go1.26.10/src/text/template/template.go (source SHA-256: `de8f0e641b4b7d8fae4cb3375fe1335280bc279d1363885db75aaa6dd5310e32`)
- go@go1.26.10/src/time/sys_unix.go (source SHA-256: `dfb95fb684453c7314fc9a31dacb6bb6de7fa3ccc1752b9bb7f2741a062911d5`)
- go@go1.26.10/src/time/sys_windows.go (source SHA-256: `2136406a64f507a80d90aaf03207f38b473629e16c6bff77811b0bb339e74310`)
- go@go1.26.10/src/time/zoneinfo.go (source SHA-256: `03bc9ba0ed7f4657121992620a83621f2ac7c5a6eb8b559d46600829e618211c`)
- go@go1.26.10/src/unicode/graphic.go (source SHA-256: `46a43107a87c88cf17a05a557e00dbf6b0f4bb53ebcd7d3035e02f1e67c428be`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/composition.go (source SHA-256: `3d1be52960f2693926472819b747646e7d2371d2bb5f53097cf9d46c913052df`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/forminfo.go (source SHA-256: `b8d7e1fd095038fee9db7ea058bea492c55d642de2178475569ee763eebe25d7`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/input.go (source SHA-256: `965b431790bb139543d71a9c497920ef7d9a15af417456a2bbd0cdb629330e8d`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/iter.go (source SHA-256: `16a40d326f9b37c5430a18894f718007728ce8324afb08d8c0a17ce7289714b0`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/normalize.go (source SHA-256: `b9e7aeff51e6cb036ff4e8c6255410e7068c112e329a137d333e639f63c4a8ab`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/readwriter.go (source SHA-256: `d600d1f6cff2536dbb49a0f6bb8c06f6fc05f2754dad2dc05f1d44f9e8e48a35`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/trie.go (source SHA-256: `d87793d558251ee8824954f0b7bc5564803e4c9d59a8c4eebb4c3c5cfbd19492`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/hash/crc32/crc32_amd64.s (source SHA-256: `a72c2ee64b5035a506864ce6ecba6823ee7260862353970c05a96cadda053861`)
- go@go1.26.10/src/runtime/cgo/gcc_setenv.c (source SHA-256: `120a9a088f87be81a95270eff8dc431fc112a1d1b803c72176c469379eb0f96f`)
- go@go1.26.10/src/runtime/rt0_windows_amd64.s (source SHA-256: `ddd23b67d4a37d15edbeca1adfb40559fafcb897eeac21fe5c37ff782fd6ec06`)
- go@go1.26.10/src/runtime/sys_windows_amd64.s (source SHA-256: `24efe0d0cc1f7743cd4bfc889fbc6dabb73409bd865c7ec7a1e86eaa01863516`)
- go@go1.26.10/src/runtime/time_windows_amd64.s (source SHA-256: `637661f93b27e0b416ba78f7358c36fd3b5b54de3684edc10ad175615e1d446d`)
- go@go1.26.10/src/sync/atomic/asm.s (source SHA-256: `0d82e9b2a1e68008568b03a6eda89769a33d141beac22d20124e89e951d452d4`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/bytes/reader.go (source SHA-256: `5b42a3100971836c10b1e0d3699b46cffd41c8be23d09e1642ee42ff795dc6ab`)
- go@go1.26.10/src/crypto/internal/fips140/aes/aes_asm.go (source SHA-256: `21dbf90094928deace73d5a7b1e472b95a91bd5c316f53eb48ad5b45d7672439`)
- go@go1.26.10/src/crypto/tls/ticket.go (source SHA-256: `61b0247979c73930db8bac363c5c6a5cf7b38906fd4be67f5a9967c8818fb1df`)
- go@go1.26.10/src/crypto/x509/pem_decrypt.go (source SHA-256: `062b4025a467daf8eb1955f4afa4ce8ed19d7234274a28edf724062eae944b6f`)
- go@go1.26.10/src/crypto/x509/root.go (source SHA-256: `813faea4e4990c9760b5fe99d8ae91c8f26c4d3b5dd3c4a6d7bef7b0f11dbae7`)
- go@go1.26.10/src/crypto/x509/root_windows.go (source SHA-256: `16921a918e7a6ed7d6db895bc9dbc12477adc0cbbc3c958db9e8f018c018af1e`)
- go@go1.26.10/src/crypto/x509/sec1.go (source SHA-256: `624dbfcfa0eb07b5860bfdff475fbc38a65dea2df8da3f0ba259aaf1a4c37dd4`)
- go@go1.26.10/src/image/jpeg/scan.go (source SHA-256: `bb475bb88f5b57a6b55e88a018244cdfae9d7c89487471050ab449f89e1bc39f`)
- go@go1.26.10/src/image/png/paeth.go (source SHA-256: `d123c1f3f08963a359d0e5adaed2526120dab1bbfc8d5e08ac59802c6748cbef`)
- go@go1.26.10/src/mime/quotedprintable/reader.go (source SHA-256: `8caa6835cbdea786381063601d021dc775751bf39c5b8d41a965f6370648227b`)
- go@go1.26.10/src/net/lookup.go (source SHA-256: `a4991e7764d27c73ffb219cfecf35df03b029911811f57caaef3f7a5c190a4e7`)
- go@go1.26.10/src/os/signal/signal.go (source SHA-256: `3e41e46229412d48f34e216a2600b58106a5f425958b12a58258913d69d65856`)
- go@go1.26.10/src/os/signal/signal_unix.go (source SHA-256: `f0f15b42aa5d147baad17a03d1417a65cdae9df997e0589b1e6e104f54d052f4`)
- go@go1.26.10/src/os/sys.go (source SHA-256: `98ecc8ed73af7c43a643a62a2a16d445cec49dacefa266dd1ef7197188b84679`)
- go@go1.26.10/src/os/user/lookup_windows.go (source SHA-256: `916d451cba64f01d92fdd2ac7e2e84d1d00b02613773ac0c72965098a5b08360`)
- go@go1.26.10/src/path/filepath/symlink.go (source SHA-256: `d4e5ae166c208c8d0fbbed9aa907208f7cc541100b8e04029a38f16f1b5da1b5`)
- go@go1.26.10/src/path/filepath/symlink_windows.go (source SHA-256: `801a4135e7dbc63dbf392fb1b4dcfcd0c8b5d81c66dfaa8aff7ea98a8f708ccf`)
- go@go1.26.10/src/reflect/makefunc.go (source SHA-256: `7b56f50cc2ca5dc2bfdd65977cf7501d049d9e66925cef0baa9e8be42689def9`)
- go@go1.26.10/src/regexp/syntax/doc.go (source SHA-256: `d61fed6f5a70cb4a0883f2c70da086d52edc80faf79ecb33a95393f28f24711e`)
- go@go1.26.10/src/runtime/compiler.go (source SHA-256: `bfe09a6aed85f21db093649c9b431f76353794d3b630b2f788859ea96a2ff6a2`)
- go@go1.26.10/src/runtime/env_posix.go (source SHA-256: `351f359fab6974c595ab72454ee246b0ebd82c3ca8bec599e33f49d6d9cefeb1`)
- go@go1.26.10/src/runtime/lfstack.go (source SHA-256: `2ca3c8998179a3212375e432618c6a7321ba17bb591a58e26f9b501e97abe50d`)
- go@go1.26.10/src/runtime/signal_unix.go (source SHA-256: `bea6344606adbcc50bfbdb6989bd3f1743ff8a7a8cd4fd66a35b9ae065f9dc53`)
- go@go1.26.10/src/runtime/vdso_linux.go (source SHA-256: `2667f1d8d2ba2b1bb8ee01c809208c04662819a7eb7e0ef8b041c18e6b7cf5f6`)
- go@go1.26.10/src/runtime/vdso_linux_amd64.go (source SHA-256: `6763e3a583e07d928cc9a54d9b1d682969883c9e38be457309748c09b1243908`)
- go@go1.26.10/src/strings/search.go (source SHA-256: `4fdef1f79f6e854789d9140b17ce2c2bbdbf0d4067a801e42607cf7877698cd1`)
- go@go1.26.10/src/sync/runtime.go (source SHA-256: `a2be85512b59e11270b099f264581343cc1589bb782e03a476102406defd5f33`)
- go@go1.26.10/src/syscall/security_windows.go (source SHA-256: `bcc99553aa61028883a744ae7e62383be648c28e199bb469e60c1770d31ea9d9`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/internal/poly1305/poly1305.go (source SHA-256: `c4b20ab7330e47ee155a775f69ae3b0e432c75d74ad005cf764228a69ddf6da6`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/internal/poly1305/sum_asm.go (source SHA-256: `ce5f94aedd0ce3349a8a8655607bd0159d354902e3298169cf9350fe235382b7`)

````text
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/math/floor_amd64.s (source SHA-256: `73cbaf1ac9f05bb7989a50fc82cce01f6ab20162ba03c8862819416cb2a270fe`)
- go@go1.26.10/src/os/signal/sig.s (source SHA-256: `0ff3cce2e20637f8abe9ac93896bbdd2d7327a15696ba4cee57274872e560c2d`)
- go@go1.26.10/src/reflect/asm_amd64.s (source SHA-256: `1a32c662e94b681985a500708c79cc33854e186f4095108b95540b0aaec538a6`)
- go@go1.26.10/src/reflect/asm_arm64.s (source SHA-256: `695304e15fbcd8f3681ae9a0a7f32955e52784677b33d6b64c6c09ee65a9a427`)

````text
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/bufio/scan.go (source SHA-256: `67ad04ac41a859798cb26e112d52ac15ec2fb11b4948293e20c0f7983d029c0b`)
- go@go1.26.10/src/crypto/elliptic/nistec.go (source SHA-256: `19452fe2c70d7b0136868829fc6466fe6ee8394c996d78bc1e7a4d27964a594b`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/gcm.go (source SHA-256: `bcb49a1b0727d616f716d1a2d6c67f8eb63d9095ab24b1424797ba08cd43427f`)
- go@go1.26.10/src/crypto/internal/fips140/rsa/pkcs1v22.go (source SHA-256: `22c489bc7f9b495d51e4f2c01c5f12e30bf1bbe289d12bd08bf611d2e3b82134`)
- go@go1.26.10/src/crypto/internal/fips140/sha512/sha512block_amd64.go (source SHA-256: `097bf048d2345d63735526793797d5c17f56bf6e583cbf1b8c24ada241705a8c`)
- go@go1.26.10/src/crypto/md5/md5block.go (source SHA-256: `fa70d9bea91b96db4525a0e9e961dfe0e967f826ebf1456940dcf63127fd90f4`)
- go@go1.26.10/src/crypto/md5/md5block_decl.go (source SHA-256: `5611f505caa831d8c41d9968ce7baf06720770bb1da252f8a95e227d8091a568`)
- go@go1.26.10/src/encoding/encoding.go (source SHA-256: `aacbbaa46c79ac0a47469a880b51b8e967e5dba80e2ea39ec0d6ec72441db08c`)
- go@go1.26.10/src/encoding/json/fold.go (source SHA-256: `5b3c91bfbf6943e04d9aeb934031331a3f08a83b10079a0a7f0def3466b2f39e`)
- go@go1.26.10/src/internal/poll/fd_mutex.go (source SHA-256: `ad10a59470d3848a745f9549afcee67aae9553e385ec41988d7e5266c2bc3902`)
- go@go1.26.10/src/internal/poll/fd_poll_runtime.go (source SHA-256: `88fda43019c5dcf42545effdebd1747aaea377560dd041a35df8499538ee07c5`)
- go@go1.26.10/src/internal/poll/sock_cloexec.go (source SHA-256: `3ae8a5856c703f68eae5343f575727ed6123310c8d4bfc0ebf1814efa8028cfd`)
- go@go1.26.10/src/internal/poll/sys_cloexec.go (source SHA-256: `501d0ecc3dc29e8b0d36d01caf2deb495c4545126dce48088ec00d5dc9b61fdf`)
- go@go1.26.10/src/internal/singleflight/singleflight.go (source SHA-256: `582ef7cfe2bdac797afea4706dc605a80234db69870a80ef492df7c9e1ea1df4`)
- go@go1.26.10/src/net/sock_cloexec.go (source SHA-256: `07f29fbd70688bfd54b0ff5c4c46af7995a303135136705157fd18cc8c80776e`)
- go@go1.26.10/src/net/sys_cloexec.go (source SHA-256: `9fbcd693976477b5a7630f8b46f2bbc556c8feadb4e580a8c10b73d57fad46f0`)
- go@go1.26.10/src/regexp/syntax/perl_groups.go (source SHA-256: `d060a0645b60b210eaec7c3dba6d554dc297c14fe71cdbcc9c25b4fc7562c3a4`)
- go@go1.26.10/src/runtime/netpoll.go (source SHA-256: `f7c606d3fa0a3b9499fd5f42b81918a98f0b99a6029fca2010cabc710fc9ea65`)
- go@go1.26.10/src/runtime/netpoll_epoll.go (source SHA-256: `e4b6d6acbb1cc36a552cd64f7fcbc7f21962f2edf88ca6e75193f673c4e4ae92`)
- go@go1.26.10/src/runtime/netpoll_kqueue.go (source SHA-256: `7ea25693196718940e852554427d6b0d96d2964899e91e1bd3852d3f9267a072`)
- go@go1.26.10/src/runtime/netpoll_windows.go (source SHA-256: `9a0727cc4b9fb93de904123083f4f9fe05af1f9b14a67218a3cb812eb1f9ce62`)
- go@go1.26.10/src/runtime/signal_amd64.go (source SHA-256: `131fc8cdf581e82194de1afa3a7d3d982f056f35e30fcd23d412bae669d0b90f`)
- go@go1.26.10/src/runtime/signal_darwin_amd64.go (source SHA-256: `25949c7753f4030764453ae291cc3770f2b086b02708cc846daa0f4e780f779f`)
- go@go1.26.10/src/runtime/signal_linux_amd64.go (source SHA-256: `ff548a1eff11973181b47df60959ec41a9ea3613c58fbb44b471368a58f92ec0`)
- go@go1.26.10/src/runtime/stack.go (source SHA-256: `0225942ade2eed1e4ec5a0e813c66b3ed06c142c33c987b2fd9e031a0b871d3d`)
- go@go1.26.10/src/runtime/sys_arm64.go (source SHA-256: `74b15ee351dcc6dbe624607ee8e0b26e76b352b3132f40f16358336af0461c67`)
- go@go1.26.10/src/runtime/sys_x86.go (source SHA-256: `f9c1c32030a7fb02db5abced4c69f250bf8bd341ac6a24b71a781908bf20c42c`)
- go@go1.26.10/src/strconv/isprint.go (source SHA-256: `24114557c2b3cad7b80c8e71683831bfbe1c72eb99e8fd28a3a48eaf77705fb4`)
- go@go1.26.10/src/sync/pool.go (source SHA-256: `d03a971c04ba5be782e997ba2c20e24d313f3540595b0a51b52f0081143a29c2`)
- go@go1.26.10/src/time/zoneinfo_abbrs_windows.go (source SHA-256: `fdabcbf35b6e9613c7e12bc44a10701cabf08c1760aef7be07f58609dfb89ead`)
- go@go1.26.10/src/vendor/golang.org/x/text/transform/transform.go (source SHA-256: `7716301e210f42b764ec78f4adffb05075f1d5103446a40b8a38043f5a7abd88`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/norm/transform.go (source SHA-256: `6f8014595643e2acae76d47ed6abe8ace969a0c9070dbfa5a89c86bebcec812d`)

````text
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/bigmod/nat_arm64.s (source SHA-256: `e48ea304e2433cec8d480cac7d326cef2556eef5dc0852c0228c3c8280b55d53`)
- go@go1.26.10/src/runtime/funcdata.h (source SHA-256: `bd73b7c27a4d3f8560a5a9dd1a5da54e0e48572311a744cf7843e709b7f2a015`)
- go@go1.26.10/src/runtime/textflag.h (source SHA-256: `a9d761d03aecdcc117691b6a55eacd43eb7d47d9ad9ca7510b8669ac403b380e`)

````text
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/context/context.go (source SHA-256: `765f99f7cc0710b0117ffa183e432524f429c21496fcbad9b63fb9886fdbaac7`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/hashes.go (source SHA-256: `ef12ea1350256c1e0201b4464eff57d949bc6b5daeb6113e48924877d1123a04`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/keccakf.go (source SHA-256: `4bb1be9fa0904272790ad452bafb18d4fb8adfb14075555692a3415ca7391c4e`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/sha3.go (source SHA-256: `a6b82ad004d074be55eb27dd37b76521f3c4dd1c7c7e8dc54de3047385c0889b`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/shake.go (source SHA-256: `ed11327ccbebefbcefd4edcd89e52554057d89d5f50f9cb142ff7daff78fd2fa`)
- go@go1.26.10/src/crypto/internal/sysrand/rand_getrandom.go (source SHA-256: `5597f30418c230c17a3ffaf88137779446291ddd6dfb038045e827a9a2a289ac`)
- go@go1.26.10/src/encoding/gob/dec_helpers.go (source SHA-256: `c42910246c6cda2eb3e884b908d4b2f541d220232669ba0821c4afcf44852899`)
- go@go1.26.10/src/encoding/gob/enc_helpers.go (source SHA-256: `edfb90e367cc4708e9c70a17091ccf0063355f4120fe7faa62f076cf090a3ec9`)
- go@go1.26.10/src/internal/goarch/goarch_amd64.go (source SHA-256: `a6cbffd95eb55c437ecec11276cc4d3b8dd30e99a45f0e225d660b117a9e88ae`)
- go@go1.26.10/src/internal/goarch/goarch_arm64.go (source SHA-256: `b5fcb79b0b2b907abd508a1736d210452945629a7046d59c1bd5f01ba46d423e`)
- go@go1.26.10/src/internal/runtime/sys/consts.go (source SHA-256: `c3ae3566bde7c7eb00d06996ce28080c6e2ef2d1d6517538e9d27bbd07a3474d`)
- go@go1.26.10/src/internal/runtime/sys/nih.go (source SHA-256: `62b649628d68a23ecd28ac18c01de4a7117b0edbec150cb33bc99abdd7d01ab9`)
- go@go1.26.10/src/internal/syscall/unix/getrandom_linux.go (source SHA-256: `9cf2c1d1ce02e08300dc45f44d1472754117350a6a135a9da755783f02c17c52`)
- go@go1.26.10/src/internal/syscall/unix/sysnum_linux_amd64.go (source SHA-256: `383d4857ec4f538630e4f781dba650ec48baacfd101456814a1b533b41201c8d`)
- go@go1.26.10/src/internal/syscall/windows/syscall_windows.go (source SHA-256: `bd17d11ec16f619351a85040109b5324a4d7ca1a168cacd398002acfa5fd7611`)
- go@go1.26.10/src/math/big/float.go (source SHA-256: `65fdfd0817868a1ef6d1123eca037f7a20ebac00135c688197d02d51f76820c0`)
- go@go1.26.10/src/net/http/httputil/httputil.go (source SHA-256: `4697e1776b36edd928fdbd2abd54d59aae4b78e0a99921216d2eedebf7553e7e`)
- go@go1.26.10/src/os/sticky_bsd.go (source SHA-256: `1f825ebbcc4409930c7748398dabca723ed2234328325d4dc466f9b751b5dc22`)
- go@go1.26.10/src/os/sticky_notbsd.go (source SHA-256: `559270cac07833f1f3c38495f2aad352ed83d68b196c32d18351545c7152ed84`)
- go@go1.26.10/src/os/sys_unix.go (source SHA-256: `f339da7dca4c290723a7ab94babc1b4dea401cb1be043aa8a5470296b2dc13ec`)
- go@go1.26.10/src/path/filepath/symlink_unix.go (source SHA-256: `9c8178a2b185bc24de24ebab3de927f3ab7a29b971f91c269817026f5aab091a`)
- go@go1.26.10/src/regexp/onepass.go (source SHA-256: `fd0e29803438c6ea3021c78227d4f69de163f693e82c155dc62d721eaacbd24b`)
- go@go1.26.10/src/runtime/alg.go (source SHA-256: `c2bfd26efce12e9a8783ecb88362738a0a52bda789cce1a3f5019575f076e9cf`)
- go@go1.26.10/src/runtime/auxv_none.go (source SHA-256: `d04fe172c21ece9ebed64be339f745e386c3ed3d21727323f0c2ff1a80b7f56c`)
- go@go1.26.10/src/runtime/cgo.go (source SHA-256: `8e2ba362ad55f640c319a8171d95805a6fa41824e6ced0bd0b1eae6291b82aab`)
- go@go1.26.10/src/runtime/chan.go (source SHA-256: `65b95334e2688452efaa5ebf0dde2b94052d383ed0f281658a9a8bc20a948838`)
- go@go1.26.10/src/runtime/cputicks.go (source SHA-256: `661b4602a411e850016ea45ebdab6453ce8de5d479ebc5f67758062aa8da1517`)
- go@go1.26.10/src/runtime/hash64.go (source SHA-256: `1ea3bf10946f8b0f210f9002146cbb72bed20f8f9bb40658575037bfa95e20e2`)
- go@go1.26.10/src/runtime/heapdump.go (source SHA-256: `a1d63344ec0dd6e1e9804137f5eda1122e265769996e9b565bd627fc6d7b8c79`)
- go@go1.26.10/src/runtime/iface.go (source SHA-256: `dc54b64c1bb6ac5c0e2cf765aee4b9465f5be86a472227af14c0a75c085eece5`)
- go@go1.26.10/src/runtime/malloc.go (source SHA-256: `67e74184452cc6cb5bedee04fbbd6f395822e1cffbb4f91ad4a484796d7f021f`)
- go@go1.26.10/src/runtime/map.go (source SHA-256: `f82e6e3c3a900049c9b85154f8a3e89085fec17697147b6003c7ead0a71dc763`)
- go@go1.26.10/src/runtime/os_darwin_arm64.go (source SHA-256: `6559a41847789afb846e8bafecacde532ee420420913d0be260effc2d3a6be62`)
- go@go1.26.10/src/runtime/os_linux_noauxv.go (source SHA-256: `7c2485a784ea54d6a6d1af5644f0ee241b8faad42d38f46f0e40c512a6d197e2`)
- go@go1.26.10/src/runtime/panic.go (source SHA-256: `4951c11983710ee957881d230290aa8cb130a873d83c55c86a0e4c7ecd35d190`)
- go@go1.26.10/src/runtime/pprof/protobuf.go (source SHA-256: `6e6fc4636549fa3e978b285df8a6c41176a9ca4191947e4ab3898ac9b21a6c57`)
- go@go1.26.10/src/runtime/proc.go (source SHA-256: `b0ea864c1dd251cd9146644863ae9126ca57761dff05862dc80ab0799925fec6`)
- go@go1.26.10/src/runtime/race0.go (source SHA-256: `40cacda94f0feb4ffb47516bd5d88ac4c7bf988bad78ba375c749900c743ed96`)
- go@go1.26.10/src/runtime/rdebug.go (source SHA-256: `39a5f8ec7012b69997d9e25c6a7f303b128f821c155ac2f3bd5e04b7a45ef796`)
- go@go1.26.10/src/runtime/signal_arm64.go (source SHA-256: `04142dad92a4fa31f79095033cd023ddf5f256b21824eb1f25e6a9600da23331`)
- go@go1.26.10/src/runtime/string.go (source SHA-256: `ba5aca9b7febcacdbe0e49c308768d79e09da1b2ddd4184e973d9f0c33b95c75`)
- go@go1.26.10/src/runtime/stubs.go (source SHA-256: `0572c8b87cdbd975e8fddc5e8e84240d6397445589641035f0ce95e6676ee90e`)
- go@go1.26.10/src/runtime/stubs2.go (source SHA-256: `afd9900bba6bd678c9217918847144cf5fed035b1035081756770bf092fcdc5a`)
- go@go1.26.10/src/runtime/symtab.go (source SHA-256: `556bdcf7743dc7be1c7c6bef9f95b968879c152696a4fa96afb6b6c48e2f45e9`)
- go@go1.26.10/src/runtime/syscall_windows.go (source SHA-256: `e7846a3244ba617530b5287585e35780e575ca10795a1dd90e8cbde039540805`)
- go@go1.26.10/src/runtime/tagptr_64bit.go (source SHA-256: `8788f7013736b8a86e9014e22362af25fed406f381febae17a0f6ce41024e322`)
- go@go1.26.10/src/sync/atomic/value.go (source SHA-256: `e87f0ad5c37b11b69effbb2c58b34b3fd7d4d4300b9e83a63847ff3b4835b8cc`)
- go@go1.26.10/src/syscall/flock_linux.go (source SHA-256: `13d2a434013cb2233ae64034d67fdd19f9df6415b5b60530044e7227fd8b242f`)
- go@go1.26.10/src/vendor/golang.org/x/net/http2/hpack/encode.go (source SHA-256: `a9c57bec32cdd84c8372aa9622c6d6c99a25d47133952b88a3f391f3e411bd4c`)
- go@go1.26.10/src/vendor/golang.org/x/net/http2/hpack/hpack.go (source SHA-256: `468bb8569beb1271228967ef4873307b86824754b6e69aaa7e9ed4849307a519`)
- go@go1.26.10/src/vendor/golang.org/x/net/http2/hpack/huffman.go (source SHA-256: `3d41693bd77da34093c1d334842d8cbeffb3a2f3e0fa5cf0af5e6e9e65b19be1`)
- go@go1.26.10/src/vendor/golang.org/x/net/http2/hpack/tables.go (source SHA-256: `6e8ce01dfb4d5e9e02a982581988621d2bc84b5e90708e1ffd3ef988cf8b2d5c`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/runtime/atomic/atomic_arm64.s (source SHA-256: `ada1ea4fe09a57a1cfdc2ded40c6b586c3f8a8d281fbab6dc0f950cd65a7ff19`)
- go@go1.26.10/src/runtime/asm.s (source SHA-256: `cd39f9e3a6d7ed158abadedaffcb7eebbcb1bccf5dcca5f05972f8f4de7f79f7`)
- go@go1.26.10/src/runtime/atomic_arm64.s (source SHA-256: `ce8189b02b74e99eea0220cdf1208b65e3c402e4be4c8d8ad9d6078fe6d622cf`)
- go@go1.26.10/src/runtime/cgo/gcc_darwin_arm64.c (source SHA-256: `693b573a1101ab9e3842977010ce04a8d318110fda9b33b9fe36f945b23949e7`)
- go@go1.26.10/src/runtime/cgo/gcc_fatalf.c (source SHA-256: `d407340597efac27ac1f308b5cbf11f2381f25cd944e6fd609319054d8c70fd2`)
- go@go1.26.10/src/runtime/go_tls.h (source SHA-256: `253985b0412d66fd6657b0fcfa54eef45d036c66a7cd00b0ac37a3da2166a4e1`)
- go@go1.26.10/src/runtime/memclr_amd64.s (source SHA-256: `826fb9aef6f97dda8e5a994e21ec692479712f69fefcf1469b756bdff29739ef`)
- go@go1.26.10/src/runtime/memclr_arm64.s (source SHA-256: `05ffc51694093eabf1ae8e501e8ff98cc9ef447fe00914447c9a49a902a75ba7`)
- go@go1.26.10/src/runtime/memmove_arm64.s (source SHA-256: `b6af6c6b8046c943bdfa82bca55bb418c9a6d2fd2c3ee438d95b6f493e7432f6`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/gcm_asm.go (source SHA-256: `2d9f49bfa9cb47b6a17ff03c6c0d8708611270e6f35b3981868d026253819031`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p256_asm.go (source SHA-256: `e98325f1f57b94dba715667f70b4eff9482a7b06ae78acae2c18785eca00d720`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/sha3_amd64.go (source SHA-256: `1c0f5d80c34c5f077475f7e0f0f46a207d4d93a18e4e52915b405a2ac5c3c80a`)
- go@go1.26.10/src/crypto/x509/root_linux.go (source SHA-256: `c9a83237f919f40fe0c54434fa3a709b71adebb7c202eb5dc36fa32bbd03e18c`)
- go@go1.26.10/src/image/internal/imageutil/imageutil.go (source SHA-256: `115214338cee7d45d39b15a63e6219ca17daea17d59b36866244bd25d6aec0ed`)
- go@go1.26.10/src/internal/goos/goos.go (source SHA-256: `2594d66acfc7136883b255809773c23c44f67d22c9fd7ca849a653991081f27e`)
- go@go1.26.10/src/internal/poll/hook_cloexec.go (source SHA-256: `3571dcf0ef879722dba04b11e8a0c72f20b5f7fddc96207a1c23e7e0420033c8`)
- go@go1.26.10/src/internal/race/doc.go (source SHA-256: `4e6670414ce0c94fffd41e26536671f0122f8ac54bdfe65277de014a8d1fd4dd`)
- go@go1.26.10/src/internal/race/norace.go (source SHA-256: `d195ef03f406d7a3e636709cccf6fc55e0fc99aa4a6a5291800364b4c308d828`)
- go@go1.26.10/src/internal/runtime/atomic/atomic_arm64.go (source SHA-256: `bfd104f7456854fd1028e02ac880380b284a639014be934357f1e55a4ed069db`)
- go@go1.26.10/src/internal/runtime/atomic/stubs.go (source SHA-256: `5dd8e4ce65dc4f6359ab8f18873dd1903f06d95a3fdcc5e184ffc2624a300e75`)
- go@go1.26.10/src/internal/runtime/sys/sys.go (source SHA-256: `15c88e2f32c3e0cd3a04a273c89372b2382430a6093741dbd03122f7f38eadf2`)
- go@go1.26.10/src/internal/syscall/windows/registry/key.go (source SHA-256: `226a71e70db92c0ba30bfb46c703366d25993957b01313c0953b64c20e9f4e1a`)
- go@go1.26.10/src/internal/syscall/windows/registry/syscall.go (source SHA-256: `904e467202debea51004211df79429652e1747affc16130084f4b842a52e0670`)
- go@go1.26.10/src/internal/syscall/windows/registry/value.go (source SHA-256: `0ce0bf86a9d5cc78dfb72254a728f0d1daf3972bca9e52719849e827d6ff6da3`)
- go@go1.26.10/src/math/big/decimal.go (source SHA-256: `bd404e9f54251bac011f1220452e708a8a801837cbdfc759766bd7b9d374f3aa`)
- go@go1.26.10/src/math/big/floatconv.go (source SHA-256: `a4bf1a47136a412b957a57eac59bce9509e600ac76a39f0de37a7703a2b47af7`)
- go@go1.26.10/src/math/big/floatmarsh.go (source SHA-256: `c60b612f19163a05116bd24bb4f30c7b869a60bcdf628a6b8a24a91a09339d87`)
- go@go1.26.10/src/math/big/ftoa.go (source SHA-256: `f47b86fadc24de328ecd7fb92677311b1deb37637bdaadd78efeffd3f54d9240`)
- go@go1.26.10/src/math/big/intconv.go (source SHA-256: `facceb8aaa8f8741637847f7d5d418980dfaa1f485bc6e74d5d11dc2182468d9`)
- go@go1.26.10/src/math/big/intmarsh.go (source SHA-256: `159c40dbd413af19e3f5ea6fc80668384bbe256334cc4cb937f22d2a902cde87`)
- go@go1.26.10/src/math/big/natconv.go (source SHA-256: `a5ba99ff1e8e0fee8ff50958a0360ad0516eee7b7208d965e10f9b99e0289820`)
- go@go1.26.10/src/math/big/ratconv.go (source SHA-256: `0eb73fc2f2824c7c3a3ad878fa132f4abccf836d301ec85acde27ced0a0e2500`)
- go@go1.26.10/src/math/big/ratmarsh.go (source SHA-256: `a7d6e4709e9a9d99b4b4061e05a0e5b78a3b317b1b663322c6eafd45ec406482`)
- go@go1.26.10/src/mime/encodedword.go (source SHA-256: `394cfb382ac4b68806b067757f6fa77ebb009392eab1d9d7a12db8e03fda838a`)
- go@go1.26.10/src/mime/quotedprintable/writer.go (source SHA-256: `c180dbcf62cfe6c686698b698883891703bdf93da7d27c722dac85e628683114`)
- go@go1.26.10/src/net/addrselect.go (source SHA-256: `3cf4ce47881e417d3e4c7e459125e3445b22714fc19c4b40cf776ab1bd835189`)
- go@go1.26.10/src/net/cgo_resnew.go (source SHA-256: `effefda66ef7cd9ddbbe2afd0f0866afd383c1c8a7d97abfb154bd6110091f9f`)
- go@go1.26.10/src/net/cgo_socknew.go (source SHA-256: `7cf6c80df52eb48299bdf38c80d45bcba058a34a11e863a19b33557fd276d92f`)
- go@go1.26.10/src/net/conf.go (source SHA-256: `c2ce5b531decce68f10df916bc353a74c564a388156cb4f2d2ed2b6800571f5a`)
- go@go1.26.10/src/net/file.go (source SHA-256: `50c0d5f95b2e27ea31dad71d2b6f68f6265d8e0d64d75d192e08b5ae322f21e0`)
- go@go1.26.10/src/net/hook.go (source SHA-256: `d4768323362e325c9500b53f6df5502dc89278b960862d43f451a61a6a1cf2c3`)
- go@go1.26.10/src/net/hook_unix.go (source SHA-256: `dc125f4a744316c9b34826e4b67ef43ea1608430cd7b79e5fa06ad3acea35dc7`)
- go@go1.26.10/src/net/hook_windows.go (source SHA-256: `54825a0db3907f4286d0ea141b244c1168240ca6e9431a376b55d0f4e521be1e`)
- go@go1.26.10/src/net/http/method.go (source SHA-256: `e0f44099c09a6245c475585347eed68c738e730008700964685a9fc7afe5963b`)
- go@go1.26.10/src/net/nss.go (source SHA-256: `b93ae402e153706ab688ac75cffc19aa8a097b964d937725d24ef9f7fac791a5`)
- go@go1.26.10/src/os/exec/exec_unix.go (source SHA-256: `b7606c80ae631012d0d1fa17acdbfd5ec36daf2c1810138d23df923f4665163a`)
- go@go1.26.10/src/os/signal/doc.go (source SHA-256: `1336acd92452015fc2db61a2c6376828109344fdc8396471f1271b164c7fc03a`)
- go@go1.26.10/src/regexp/backtrack.go (source SHA-256: `ac5a1e4cf49b1c7f4372c0a7c2e8d156100c17d3a114fcf9539f55d223d06092`)
- go@go1.26.10/src/runtime/cgo/mmap.go (source SHA-256: `c1c248df38c648d3880ff2fba9dd98a9c032c9bf98fe714f76eac68cce98b459`)
- go@go1.26.10/src/runtime/cgo_mmap.go (source SHA-256: `45f6998ad658f078d5b8800f62b60c407d8c637a8be69330579c82a26a3a5db2`)
- go@go1.26.10/src/runtime/cgocheck.go (source SHA-256: `0f201428e80a9b221653431811cf7e894049602c26ce1efb54a8717cb0faccc4`)
- go@go1.26.10/src/runtime/cpuflags_amd64.go (source SHA-256: `b9aa3542d35b84d5279030a603acfade2563c67b970e99d5d89ebf64742383f6`)
- go@go1.26.10/src/runtime/fastlog2.go (source SHA-256: `4e03110f81e7323d3924fda1b76e3522e0631cf3d107f5c29de8be0b124a14a2`)
- go@go1.26.10/src/runtime/mbarrier.go (source SHA-256: `a7a24f880adbf607fa8997ac7f14b6999ba466b395a99ce01baaf39e361a0122`)
- go@go1.26.10/src/runtime/msan0.go (source SHA-256: `a53ed6335e7b05d189b422b3503d2421203362754eaa7b0b73c7bab4780fc2e4`)
- go@go1.26.10/src/runtime/signal_darwin_arm64.go (source SHA-256: `7fe01f6eea07e8b4ed54cc524beb33864feec0806c1b74ab3e05c9bb6c44baa3`)
- go@go1.26.10/src/runtime/sys_nonppc64x.go (source SHA-256: `fa63b093b68b24b280878384269e4ef7c9791216991e616ab3537a1484c07947`)
- go@go1.26.10/src/strconv/doc.go (source SHA-256: `aea030aeaf555c48ae648dad6421aca2cb9df369265ab31bb1cdf956dac8ce97`)
- go@go1.26.10/src/strings/compare.go (source SHA-256: `af4e0fa7dfe63018d9808f89a4c511eeb8d773bd6f3b030f1c9592a1e0de7a85`)
- go@go1.26.10/src/syscall/syscall_darwin_arm64.go (source SHA-256: `3007eb168979860bf1b363a0916e0dc195ff0af7eca660d01885b7fb5792b941`)
- go@go1.26.10/src/text/template/option.go (source SHA-256: `fae2206d9c10cf31439ba9d616a2f8ead93f8e9b8fc12060fbf9cfbc842c0a43`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/bidi/bidi.go (source SHA-256: `b1e9360e9a026006f56914a5716b18780346418bb0b64bce23a4ff92d8147b13`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/bidi/bracket.go (source SHA-256: `5caa24f9c4e04a54d18875468f976b4b7c09c8700e88286f2fd96828f134107e`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/bidi/core.go (source SHA-256: `1d375d88eeae31f2496e081df73c0004b5c4b4a800ec6f4aea989d18fee8ec92`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/runtime/atomic/atomic_amd64.s (source SHA-256: `5b8ab1a7f5ef5aef1dc268d1a29b9cac53feaa271b995b99613506f3dc20371c`)
- go@go1.26.10/src/runtime/asm_arm64.s (source SHA-256: `c2d54a06306c90a1d6ab666101c563056578a84c9b11f50575874af9c54b2133`)
- go@go1.26.10/src/runtime/asm_ppc64x.h (source SHA-256: `443dd7c93c404f7039a48acdeb7a860ace231b3d2b45daec9ed926e577b7cf61`)
- go@go1.26.10/src/runtime/cgo/asm_arm64.s (source SHA-256: `bcc717d1b6e73e519956f722b19518e7f9bb627ec6d5f9aa6ddf21fde7f54775`)
- go@go1.26.10/src/runtime/cgo/gcc_arm64.S (source SHA-256: `d57a3bdbc4b2dca74a4a0fa834e1dbb916abc237e9a1f08855f5274bcb884514`)
- go@go1.26.10/src/runtime/cgo/gcc_libinit.c (source SHA-256: `dac91ed9f69b6762a8238215a59777f3540daf67dd93c8b1e898d01c62b69826`)
- go@go1.26.10/src/runtime/cgo/gcc_libinit_windows.c (source SHA-256: `d3fcc49070277604a633eaa3875acd65e67df343b61f829afac552b91f5fd0f9`)
- go@go1.26.10/src/runtime/cgo/gcc_mmap.c (source SHA-256: `d8ece6e6a4385bf0449ab50c51887946e39c49bd6d43b3a5f3757a7546c73c96`)
- go@go1.26.10/src/runtime/rt0_darwin_arm64.s (source SHA-256: `136bc79ef0ce9a6c490e8caddeac67836d331406acc6f89ee579d88d0ad7621c`)
- go@go1.26.10/src/runtime/sys_darwin_arm64.s (source SHA-256: `1de4c930a39379797b89833e3f2248c6d2fbc46ee8c5a5462490486756658fc4`)
- go@go1.26.10/src/runtime/tls_arm64.h (source SHA-256: `4a05646e67b5154afa82605427c1176c1b113023297ff5d9cbd9755a6c79fe57`)
- go@go1.26.10/src/runtime/tls_arm64.s (source SHA-256: `b710ccbed6a02b80cd44a6289fedc6adced85b91ca67fb7893bcfa57a99b4b0e`)
- go@go1.26.10/src/syscall/asm_darwin_arm64.s (source SHA-256: `20a0a3d681ae7801810665cb26b34d762cf82fae21e1883c5a24c7d2b2958e16`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/compress/flate/deflatefast.go (source SHA-256: `778d7457e4c48d2f6133d22e9f08c72e3299c0cb594c02d44cfc4d229cf05b08`)
- go@go1.26.10/src/compress/flate/dict_decoder.go (source SHA-256: `26ce1ac6f0a735d8e7fffab3f1052d2e2c6929e9fe740ce7c514a7f46552b28a`)
- go@go1.26.10/src/crypto/ed25519/ed25519.go (source SHA-256: `4f931fb1680d5d7a004e38660ff3d456cc4916c0224a4599f8bafe2d37d262fa`)
- go@go1.26.10/src/crypto/sha1/sha1block_amd64.go (source SHA-256: `f87d1c5838f47e50052ef6c694e366630a75c3996abb16f392d3e7a0fedf56b2`)
- go@go1.26.10/src/encoding/json/tables.go (source SHA-256: `74f7f392e1cc22890b72959647b9e5facda591a890c81f8167b532d31adc561e`)
- go@go1.26.10/src/internal/nettrace/nettrace.go (source SHA-256: `0da7a09e33953e6e690b78b2eb30749bc433696dddc24d8b0e158c17274b6b96`)
- go@go1.26.10/src/internal/poll/writev.go (source SHA-256: `e7252e57fa08d577f354ffc703d1e7d640975e09b8c79de9ef682fce93fefad5`)
- go@go1.26.10/src/internal/reflectlite/swapper.go (source SHA-256: `9e94a5d672db8bf6e0b3a4dba575e0345e069a0d2fefafc7c98ac9a93e126f64`)
- go@go1.26.10/src/internal/routebsd/address.go (source SHA-256: `329d69879b3a96924de04eb86a2456484154246323c8c9714f0019273bb85703`)
- go@go1.26.10/src/internal/routebsd/binary.go (source SHA-256: `f8d896ed8b25856e40d82eada2bbf7ff521bdebab555de481fd97f28fb13bc9e`)
- go@go1.26.10/src/internal/routebsd/interface.go (source SHA-256: `9546725f127a70af93052c2bb70f37c7581521146ab1df502d85be39ed44e189`)
- go@go1.26.10/src/internal/routebsd/interface_classic.go (source SHA-256: `f76f262241ec54d7fc293c288d6cee2750d09afb89252f7a6285f55f496f8610`)
- go@go1.26.10/src/internal/routebsd/interface_multicast.go (source SHA-256: `f707e13b6fb55cfe09cf820ee3d11718ffb518e35a21aa705912d36a68a1d124`)
- go@go1.26.10/src/internal/routebsd/message.go (source SHA-256: `fdf8c07f8968dfeebdd70f8e25025f57a04099a0564aa5b730483c0ec5719366`)
- go@go1.26.10/src/internal/routebsd/route.go (source SHA-256: `34d5f09f2fb078369e8d97ccd98ce91a620a37d20a0a247c4af8d574d39fa513`)
- go@go1.26.10/src/internal/routebsd/sys.go (source SHA-256: `fe43a2559d14c525c756e5a03d040449d4203a4461325374a43fa9478c34c799`)
- go@go1.26.10/src/internal/routebsd/sys_darwin.go (source SHA-256: `5bf3e1def9c0a7a6c8b71c110c5a8312bb5d8f3dd7e006ce5553e025a7bdc422`)
- go@go1.26.10/src/internal/runtime/sys/intrinsics.go (source SHA-256: `b7e4af48817954e8b86ff1a0c48b8c94b42e6d7c70230cf2c9ee3ad4d8f80be7`)
- go@go1.26.10/src/internal/syscall/windows/reparse_windows.go (source SHA-256: `a27f00803580d0d994bbc908ac6707e4e713e0499e5fdfcf914197e79c26471f`)
- go@go1.26.10/src/internal/syscall/windows/security_windows.go (source SHA-256: `420addd3f61110cf0d4320d326b5f2442f0148bc18f416514fc4db779773ae46`)
- go@go1.26.10/src/internal/syscall/windows/sysdll/sysdll.go (source SHA-256: `41b9bc832549766c2d1669a548612f181760effd45355c4a831f30118598d687`)
- go@go1.26.10/src/math/big/prime.go (source SHA-256: `783e8e94c38db53b28b67ffacecbdb6efb7b04f5ee50e77666a530ac10956609`)
- go@go1.26.10/src/net/http/http.go (source SHA-256: `bdb1566fa95a7667ccc9c58d5756a3730ff32d19fbdc108fd204b9179ae7a9ef`)
- go@go1.26.10/src/net/http/httptrace/trace.go (source SHA-256: `7a60b75f019f4740d210e0c11c3214a6c679f0a5ddd8865b6073790e56d3d158`)
- go@go1.26.10/src/net/port.go (source SHA-256: `1394d22f0da60f2f1432340642800b6314f1f505ceac89892c1a63879d36c23c`)
- go@go1.26.10/src/net/writev_unix.go (source SHA-256: `01c0e12624b32af49ef3f92f9f8b4c92e92522ab75c082e1ed131ba7be2573c6`)
- go@go1.26.10/src/os/dir.go (source SHA-256: `bc8540932fe12e7e5325ec8f0075eaa5014241343ea45f0f9b755f00fc672efc`)
- go@go1.26.10/src/os/executable.go (source SHA-256: `28709f736a30c43977d06a799ce86baa3dee133d11ab7583d53b07de38d41bcd`)
- go@go1.26.10/src/os/executable_darwin.go (source SHA-256: `cee8d5195663341f75bbbbba2c96dae42afa406ccf9a618e17e858b720ea30b3`)
- go@go1.26.10/src/os/executable_procfs.go (source SHA-256: `61240d3520f1294d95bf0de8e8a324367cd66b06727b319ecb55e2b3046e5e46`)
- go@go1.26.10/src/os/executable_windows.go (source SHA-256: `336118c7f02ef311ea7901d86956290079a6683a6d807a52638a8ffe2a80e2f1`)
- go@go1.26.10/src/os/stat_unix.go (source SHA-256: `fa12fa05805074520f72342648764d156193805d6ddfa663af938145769a204a`)
- go@go1.26.10/src/os/user/cgo_listgroups_unix.go (source SHA-256: `e7eaf76e32311ed493164f397f9813f6cd39c713281619ebd33d156ab5ef2048`)
- go@go1.26.10/src/os/user/getgrouplist_syscall.go (source SHA-256: `37bb7c833cd78af16e203707c879272cee98915571f44d2dc4582c5732633cdc`)
- go@go1.26.10/src/os/user/getgrouplist_unix.go (source SHA-256: `c3270991ecea12944ee46ce9633f7df85b7181dcd1e0bf1a7f2213c877e8f674`)
- go@go1.26.10/src/os/wait_unimp.go (source SHA-256: `1aaefda3a61a0d20f5f374f6181c1eced43612df07a8e82948c6124f4a46abcc`)
- go@go1.26.10/src/os/wait_waitid.go (source SHA-256: `7c2030455bb3fd966459edee00f1dc11e0034689bcae84fd44fdd155e6a0d2a1`)
- go@go1.26.10/src/reflect/swapper.go (source SHA-256: `5d8467c8475770fac34d48536e36e40d4251612af2d5a7d96f19fff1c7d4ba99`)
- go@go1.26.10/src/runtime/cgo/callbacks_traceback.go (source SHA-256: `e5046003c58a73e692b4045763b0faf7a6c955fe199a017cef9a18c9383086c3`)
- go@go1.26.10/src/runtime/cgo/sigaction.go (source SHA-256: `d55a1f20e54e2c922f5b19c2b9323f7820e70371bcf0d00ac4e2adb8c614e5cd`)
- go@go1.26.10/src/runtime/cgo_sigaction.go (source SHA-256: `cbd8a16fea9bd4e71f05b113869d0a4b203016e728e773499570623a8e6fe040`)
- go@go1.26.10/src/runtime/plugin.go (source SHA-256: `233fee9392e7bdbb1e2847b4538112eb5492331176f874e6f2f700cc6913a37d`)
- go@go1.26.10/src/runtime/pprof/label.go (source SHA-256: `887b26fce1e4ce50c76d2db6ab6881939de88e082cc088e76a01e83be46b2a28`)
- go@go1.26.10/src/runtime/pprof/proto.go (source SHA-256: `7db7c755f02551e22207987a6d5a6d25038fcc5119e9005e7be4f56fdcbc54dc`)
- go@go1.26.10/src/runtime/pprof/protomem.go (source SHA-256: `4c4dd02dec9008aeae3798b9d99cea860e2d3a8ef2530214c8035585cfbc5df8`)
- go@go1.26.10/src/runtime/utf8.go (source SHA-256: `08e09814d245742077592fe53a434a630c3067273fa88321f893627dfad81a6f`)
- go@go1.26.10/src/syscall/setuidgid_linux.go (source SHA-256: `59d764551533051659477d084963f111d02eb0ae13841b683383685eff3beb31`)
- go@go1.26.10/src/syscall/timestruct.go (source SHA-256: `4271edf5d1f895c0f94f6af4c067fc71c1197e11cf3e5b85facc563f321101fb`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20/chacha_generic.go (source SHA-256: `34403e82b1387b4402b00ce30c1364508c333f4bdbe671321690c6ebaa8d3180`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20poly1305/chacha20poly1305.go (source SHA-256: `08400b984074331e97c7e2a1d9701f14cb8f9919b3eb3e0d1bb8e3549961cac2`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20poly1305/chacha20poly1305_amd64.go (source SHA-256: `706778869957b0b0e5f14c84be1a4082ec24c661bb758d19480ae9d14cfd6d2f`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20poly1305/chacha20poly1305_generic.go (source SHA-256: `f3e04c8517e163f94ec2bb6e7b634ef239432ed7422b1d8695cfc446c8551181`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20poly1305/chacha20poly1305_noasm.go (source SHA-256: `0ff7a84348991f03aa8c9a74472b112036cb8149218aa284ed034c2ea011afc6`)
- go@go1.26.10/src/vendor/golang.org/x/net/http/httpguts/httplex.go (source SHA-256: `8cfd5819dc821d45f82fde10afff4b2e652535b0091e0559fc8ad7766051a2c2`)
- go@go1.26.10/src/vendor/golang.org/x/net/idna/idna10.0.0.go (source SHA-256: `f2232848beb271df6b3d05f6df16e10cb8c9f94ecff709a5e32dd1f2bf1e64cf`)
- go@go1.26.10/src/vendor/golang.org/x/net/idna/punycode.go (source SHA-256: `4e7adef5c3c4ff789ad439094738b5829896814e3fcd8a9a2c2c2ba40fa12a28`)
- go@go1.26.10/src/vendor/golang.org/x/net/idna/trie.go (source SHA-256: `8d2a16a56e4cab9d23a7afb49443bfe7a5ff0d74e885acc9121b4a9c598825eb`)
- go@go1.26.10/src/vendor/golang.org/x/net/idna/trie13.0.0.go (source SHA-256: `059e31556bdfff9d9b98660bae14911834a5e72475696a8b07baf0466ea17021`)
- go@go1.26.10/src/vendor/golang.org/x/text/secure/bidirule/bidirule.go (source SHA-256: `93ef75e5b5e4db48b975d5cea783612062377761c918e658a612adb3b7753f21`)
- go@go1.26.10/src/vendor/golang.org/x/text/secure/bidirule/bidirule10.0.0.go (source SHA-256: `62d8a61936d6417a9094927a243da7a56b3fd6cabdec368b8491a1bc6a7a59e4`)
- go@go1.26.10/src/vendor/golang.org/x/text/unicode/bidi/prop.go (source SHA-256: `f6391b2f69a1ae2a0ac1f0599a90b260ec350a19f73ad5d6572660100b5950bc`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/math/dim_arm64.s (source SHA-256: `92afabfd2cfbabe27eb43b89048f7a48dbfda034f519ac1d078a73e23917dee1`)
- go@go1.26.10/src/math/floor_arm64.s (source SHA-256: `447cc2e99f86c00020a14fdd23264dd5fca38ffe31754f76fdfc687e937ec6a8`)
- go@go1.26.10/src/runtime/cgo/gcc_context.c (source SHA-256: `f117768145b8a31967300073bcb4a23af306b11d1cdf2facf0c438f0760058b6`)
- go@go1.26.10/src/runtime/cgo/gcc_sigaction.c (source SHA-256: `5e58cb3ab2f7acfd43e1dbdbd5988a3fbcfae73ce759187d1d155febe556ad58`)
- go@go1.26.10/src/runtime/cgo/gcc_traceback.c (source SHA-256: `1a51727ac7da836a98eb4987c897a1e1f7c73922c0de1a74b93add8ec0ddece7`)
- go@go1.26.10/src/runtime/cgo/libcgo_unix.h (source SHA-256: `56cd53ecb0c17ac9152d1a2a807d9e64224e5f608c74ff1f0916e52d38822f54`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/boring/doc.go (source SHA-256: `085f6c0eeefbae4baa9d90a6cd081d4b5f769551dec00cb8d4ddbbbd5b9353b0`)
- go@go1.26.10/src/crypto/internal/boring/notboring.go (source SHA-256: `ca7890ca233a46d31b7d7d67d7391a97956f2b2aa66a882244ba5be98a75f88d`)
- go@go1.26.10/src/crypto/internal/boring/sig/sig.go (source SHA-256: `faed644e3071cddf01e07f209561daaa7fb6f54157ba417db60d942dbd80d523`)
- go@go1.26.10/src/crypto/internal/fips140/sha256/sha256block_amd64.go (source SHA-256: `8d760c095cfed0ad72ac5db3273403c0b5cea612747778affd2171545279a510`)
- go@go1.26.10/src/crypto/internal/fips140/sha256/sha256block_arm64.go (source SHA-256: `35524999d8d08ee83bff065f10947cb5ef079f26bc5d45514199e8f7d9f18130`)
- go@go1.26.10/src/crypto/sha1/sha1block_arm64.go (source SHA-256: `135fabf9e1537200ebede9af7a5c8ec9d335d94d1581c827aff37ed990e98aa7`)
- go@go1.26.10/src/crypto/tls/auth.go (source SHA-256: `7bdcfa1d4d1b53fb8e1306bcf8712037ac870cc793685d243cab508ebe654bef`)
- go@go1.26.10/src/hash/crc32/crc32_arm64.go (source SHA-256: `e3c97e14322a0dd605f21797bf1aff2cbfac580ed41b109cabab909d35b7f320`)
- go@go1.26.10/src/internal/cpu/cpu.go (source SHA-256: `926a074453853ebcd5a16754d716a1a2fa3af30906ae32198cf06b13a8dd911b`)
- go@go1.26.10/src/internal/cpu/cpu_arm64.go (source SHA-256: `ed4c7d9894f10ab6dc471a133400b145909b6da88055afb7b990b9a3ffd0d4fa`)
- go@go1.26.10/src/internal/cpu/cpu_x86.go (source SHA-256: `d1071bde37c203929b294beb870c7a3d6d6cb0f3ee2ae45eba3a2b15beb84a50`)
- go@go1.26.10/src/internal/poll/fd.go (source SHA-256: `39051b29c8a2097c40efe00072247b0bbe4fb2b0e330f504c524de7751e909a0`)
- go@go1.26.10/src/internal/poll/fd_unix.go (source SHA-256: `689fb78bb32b8e810fa1dd2ded475a79574f24a8f40260af2d4a3273392bc2cc`)
- go@go1.26.10/src/internal/poll/fd_windows.go (source SHA-256: `a6ed6ff45e39c6c49eb42c43bc07f2290315edcd63a2fb6e4028782ac3a8e950`)
- go@go1.26.10/src/internal/poll/hook_unix.go (source SHA-256: `4a467c97720f9d92849b4d64fd95980f8528d9ae2ac2a50cfc1309ff0759cba0`)
- go@go1.26.10/src/internal/poll/hook_windows.go (source SHA-256: `b63c32b4a5f03051f4bd78bb4440232fd929ba90be8e66cb404aa2f31dc12c5b`)
- go@go1.26.10/src/internal/poll/sockopt_unix.go (source SHA-256: `b82d6dd2e70be5d69c47eb197badef011291ccd15c29652a173c7db0bb69eae5`)
- go@go1.26.10/src/internal/syscall/windows/psapi_windows.go (source SHA-256: `4c2c2ca7233d50c85c88b9035eab3fd50382a88724520fe18b4c9fe4f16be7aa`)
- go@go1.26.10/src/internal/testlog/log.go (source SHA-256: `116a76747970a636a8888c7419486a3466385145a229cc5db7ca0b9ac7272348`)
- go@go1.26.10/src/math/big/arith_amd64.go (source SHA-256: `587e7caa42b7115215add76096b8304a6b86c9718c531b5e6d1c62003d25bbe1`)
- go@go1.26.10/src/math/big/sqrt.go (source SHA-256: `d843f7b83e671bf7783fc04a71e4713e5fc086ed609cb289fc0130585bacd243`)
- go@go1.26.10/src/math/bits/bits.go (source SHA-256: `694e0114bda396589e8c993a3d7880a0b461da23ee64c272929a3a5b69c2acf2`)
- go@go1.26.10/src/math/bits/bits_tables.go (source SHA-256: `c4cb4a05dab6d395d9d51720893a42f0a9d1cdf8a0c2383d87268e49acb73f7f`)
- go@go1.26.10/src/math/erfinv.go (source SHA-256: `becf9ccd6869b7035b79d8f4902286134b8a55e73747cca78bf787a66512e6da`)
- go@go1.26.10/src/math/exp_amd64.go (source SHA-256: `4c01fe9e09324bf28d617f81ace0ad03f3bf1cbd72930a772d18762fbffb32a5`)
- go@go1.26.10/src/net/error_posix.go (source SHA-256: `029fb7191046642c43ec9e6c430ea753189a9dec1e162cb4ea3048fcfa3928d6`)
- go@go1.26.10/src/net/rawconn.go (source SHA-256: `4d3c68411c4ede8c2f3b1061d316079e050a1e8facde2019c250a0c522d062a3`)
- go@go1.26.10/src/os/exec/exec_windows.go (source SHA-256: `2bfe7b546b7223ce72882a97012121d02beb61dceb88e1c64bceff4ee9186086`)
- go@go1.26.10/src/os/pipe2_unix.go (source SHA-256: `40080f961b8a6848640f344494f843b0612b9074943f719c6485ec305ab03ef1`)
- go@go1.26.10/src/os/stat.go (source SHA-256: `fcf9bc032dc2a30b40329299db15940192cad00c8092f256c257054a95db0ac3`)
- go@go1.26.10/src/runtime/float.go (source SHA-256: `3bbaefbb410e01afc0e82477b8d26470ce3bd1b2237776604a3fd85abc478b05`)
- go@go1.26.10/src/runtime/mwbbuf.go (source SHA-256: `f738c769f502a9073bcacb4c4db473ebb44dcbe91453135c4c21a26356288794`)
- go@go1.26.10/src/runtime/nonwindows_stub.go (source SHA-256: `e7e55cd09500be799d09d7e21338f6fd18840c813b5e753e39b9f4371a84cbc3`)
- go@go1.26.10/src/runtime/pprof/elf.go (source SHA-256: `43a1284bfc363574cb0843d2d632ea86dff94fdf5388070f3d32c1fc6329f8d6`)
- go@go1.26.10/src/runtime/pprof/map.go (source SHA-256: `ee129384558f549e6443ea790a89a00766e3c984ef3dbf457a11d16d076601b2`)
- go@go1.26.10/src/runtime/pprof/runtime.go (source SHA-256: `1681e97e8aa017590f317e1992fc6225b4df6dfe9ddf29e584c50c6adac0accb`)
- go@go1.26.10/src/runtime/profbuf.go (source SHA-256: `aad42f919f9fca2a5685fc2989e56cb2a71beeab1690d6e7fcece9aea953a7db`)
- go@go1.26.10/src/runtime/proflabel.go (source SHA-256: `dd18006615fac877fb74bf06e2003833c2222e447a752ad25983a35ba3d09426`)
- go@go1.26.10/src/runtime/runtime_boring.go (source SHA-256: `b66676b383e1baaf34b9482c50e2af4e6b9fe3178af6fff638368dbd6547e0c5`)
- go@go1.26.10/src/runtime/rwmutex.go (source SHA-256: `7ab5f917f2d92166e7a2a1db37fd1ad59724ebe9cb2763fa55aa16bc1e78c8fe`)
- go@go1.26.10/src/runtime/stubs_linux.go (source SHA-256: `a93986261d9225c7c06c813037f0cc81209c8851631a74796b6e3854d7f21415`)
- go@go1.26.10/src/runtime/stubs_nonlinux.go (source SHA-256: `84eae536719177fa8bc75e08cac48f2856cfc1a7330bfaa1323c59273d88f54f`)
- go@go1.26.10/src/runtime/timeasm.go (source SHA-256: `5c060a9c53fedf7c36023bbac52589626972ecad66b43a9672f8985eb43a7344`)
- go@go1.26.10/src/runtime/timestub.go (source SHA-256: `0f4200ad6f00e786bfab7251705bc87d658ba6bf2645f5b7decb413fb0887ca8`)
- go@go1.26.10/src/sort/slice.go (source SHA-256: `842b1095d5a97a239e0d88726de1e7a7bbcf5bb7417e2d6a4dd7706dd6645ae9`)
- go@go1.26.10/src/strings/builder.go (source SHA-256: `02370c56ffbd50bbe8e8ab02f8cc0a0ecf934f02981aeb9a85b025daffdb428b`)
- go@go1.26.10/src/syscall/forkpipe2.go (source SHA-256: `8178af2b2692f288af9b08859acf52266bf8ac4e2b7fdc6cab0786d3dd50aca1`)
- go@go1.26.10/src/syscall/net.go (source SHA-256: `96c7ea22567bb7a64b294ee89f91801f116b336beeabc341d016b37cf840058d`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/cryptobyte/asn1.go (source SHA-256: `a7b568d1a846cfccee16f7847c5a0dcc5326c3c36dd299cc69956d9c1ef164bc`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/cryptobyte/asn1/asn1.go (source SHA-256: `306307e55077ec60b0c90dd90e1a7954805349ce2fd68d66f08a660524653f80`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/cryptobyte/builder.go (source SHA-256: `63e388880001b25f0b8bceb2b3702699146f3691f35dd5284a7651aff93fcf1c`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/cryptobyte/string.go (source SHA-256: `66485c8123b338056d97077938776955fe6c080e05c2c77e36501946113d68e9`)
- go@go1.26.10/src/vendor/golang.org/x/net/http/httpproxy/proxy.go (source SHA-256: `e1cde07b566be558b32a4a54cd1e8d9c36ce7d0579706c9495537154f8c1041e`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/boring/goboringcrypto.h (source SHA-256: `d02b93dd188657550805418a74306addb0a96f1fc7a3346e91dc9cb9533610c0`)
- go@go1.26.10/src/crypto/internal/boring/sig/sig_amd64.s (source SHA-256: `0fab1084bb898eec15920788d6650434f5b49bdb354c21de57a2f586979e2435`)
- go@go1.26.10/src/crypto/internal/boring/sig/sig_other.s (source SHA-256: `d221eed848d1de42ea6cce16ff460cc84be9f56f8d88b42d1ac084e51c0a7f47`)
- go@go1.26.10/src/crypto/internal/fips140/aes/aes_arm64.s (source SHA-256: `11c393df5db5fdd0d93f880ff66ab796200e7bb3228902ed0269a76250d83277`)
- go@go1.26.10/src/crypto/internal/fips140/sha256/sha256block_arm64.s (source SHA-256: `0a6ec8521fbe276ae32db5aa6a67b10f2ee465df43ba2109e2d10faff8e34ef8`)
- go@go1.26.10/src/crypto/sha1/sha1block_arm64.s (source SHA-256: `4d91b79721d14d956c8a9b0fe0b343d8edde2ff955f00b38dd91d4df46d95010`)
- go@go1.26.10/src/hash/crc32/crc32_arm64.s (source SHA-256: `aedb6699f5f65086add1e0b86c8824fcf5f29cfb08c390c19b55441b01bac1b5`)
- go@go1.26.10/src/internal/cpu/cpu_x86.s (source SHA-256: `ee4dca9916b29475e600d85f71f0dfd134496d153c77d20a643b139cc6e18fed`)
- go@go1.26.10/src/math/exp_arm64.s (source SHA-256: `9fc7ff6a566378e0b5074c55a813468cc9c88cc37d6060c414288e8198edfbc1`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/alias/alias.go (source SHA-256: `e73a5572bfd89d59cf941c0a98a4442170146c543e630b9a1063808f96ccae2a`)
- go@go1.26.10/src/crypto/internal/fips140/subtle/xor_asm.go (source SHA-256: `c7d260f6427a83bdc893fceb97998c93afe6f0b76f196c64c67b345be65061dc`)
- go@go1.26.10/src/crypto/internal/randutil/randutil.go (source SHA-256: `1e57f1405462b10699b9f02c08995de89b9744cc5b1c5e3754834f6733740c31`)
- go@go1.26.10/src/crypto/tls/handshake_client_tls13.go (source SHA-256: `7c0d21ca7805e1256b800858fc3d35850c62575a795ee61b912be964cee66e4c`)
- go@go1.26.10/src/crypto/tls/handshake_server_tls13.go (source SHA-256: `54cce7be02c286d397b9af1ded5a24ad356962ad5fbfb37bf8f93ded9c32aee0`)
- go@go1.26.10/src/crypto/tls/key_schedule.go (source SHA-256: `561561d63d99794ca7e2cf91890dd60328a5140764465abe4bf040c544c2d037`)
- go@go1.26.10/src/errors/wrap.go (source SHA-256: `098116636610ae87dd80e49c4dbe6a2b6c37918d7e1af35ded005be621d68a40`)
- go@go1.26.10/src/fmt/errors.go (source SHA-256: `01fd304868493a3aa1198cc79a56274930f687353e9405211ca063c1b8195ea6`)
- go@go1.26.10/src/internal/bytealg/bytealg.go (source SHA-256: `d13371f1998de1505505052cfdf98f27e7bd8840ce5aa4736235ec3f629a9ba1`)
- go@go1.26.10/src/internal/bytealg/compare_native.go (source SHA-256: `a5de9fede104ba03d3c627fbb837ff870896d8a730359056c5e870147ea47023`)
- go@go1.26.10/src/internal/bytealg/count_native.go (source SHA-256: `d3eff88a61aafe415c85c0e9856bf1a22a3a0806ac116569cbc3451e8e01826f`)
- go@go1.26.10/src/internal/bytealg/equal_native.go (source SHA-256: `1db4aa0b4f804a208dbe4e8894fe05949d72e9b4492ef9fcac517e00571fe02f`)
- go@go1.26.10/src/internal/bytealg/index_amd64.go (source SHA-256: `77d0b48b8915657b09447ba36b163134a0e3fd08a776668dbf2dadfda6915bb4`)
- go@go1.26.10/src/internal/bytealg/index_arm64.go (source SHA-256: `6822f757f21c5e06fb0deb772ecc22e6b9ac27799eb250b729267ff4624a11db`)
- go@go1.26.10/src/internal/bytealg/index_native.go (source SHA-256: `885de0c47293c3e4d7083972be04e365e66cc4e9411a59393b2b835fba3e3105`)
- go@go1.26.10/src/internal/bytealg/indexbyte_native.go (source SHA-256: `3f9ee9987e6f71d2e7729feece762016488547d3762db2aa20030fe449898331`)
- go@go1.26.10/src/internal/fmtsort/sort.go (source SHA-256: `46d972102d8d6900ce8d1eb24be5c8ee92f6a0a8c78aed946292ecb50ff6a81b`)
- go@go1.26.10/src/internal/poll/fd_fsync_darwin.go (source SHA-256: `b225fc3aa1509753d492c6c571d477c8d91d12e817028e5a891492433a65286c`)
- go@go1.26.10/src/internal/poll/fd_fsync_posix.go (source SHA-256: `5e19d9488df676a69838524881b9ba7689b40df87016b16100b33efa131ffdaa`)
- go@go1.26.10/src/internal/poll/fd_fsync_windows.go (source SHA-256: `cd0d51c46f328868f3ccfec0009a10add5089bfb54b6f8917ca4c287aa7ffa78`)
- go@go1.26.10/src/internal/poll/fd_opendir_darwin.go (source SHA-256: `97daa47ed38dd9d33c7f8566d1b57215cdbbff3e9e2fd3a04377a7a4d5e7a5d7`)
- go@go1.26.10/src/internal/poll/fd_writev_libc.go (source SHA-256: `e4c1301b16afa9481811fb151c97cbf2128783a6b8399e2efa508be6636d7107`)
- go@go1.26.10/src/internal/poll/fd_writev_unix.go (source SHA-256: `09196eb54cc5e20e5c8cd07a4641cae561ae10b465522c60404cd5972b598aaa`)
- go@go1.26.10/src/internal/poll/splice_linux.go (source SHA-256: `969eec3b4a59f5ab1059eeabec012e0387c14b44c50a78a9b551198e9b54c93e`)
- go@go1.26.10/src/internal/runtime/math/math.go (source SHA-256: `d772845f24651c01b98da2e2469a3017413040ddb3012d73e12823a5c130a841`)
- go@go1.26.10/src/internal/syscall/unix/at.go (source SHA-256: `6b24496de4a9b1586bcb88b6acddcccc80777f302bd390288ebe178ac8064705`)
- go@go1.26.10/src/internal/syscall/unix/at_fstatat.go (source SHA-256: `5e70a71173b2e771d34d17707de3fd14ae415dfafb318d5a8641144dd851ad6a`)
- go@go1.26.10/src/internal/syscall/unix/at_libc2.go (source SHA-256: `27303a64a0a49c6b4d06e70cebbd29feca062185d7e1eb58edce1ed8ae84acec`)
- go@go1.26.10/src/internal/syscall/unix/at_sysnum_darwin.go (source SHA-256: `8981473c4f2d235ff5f7ac6114eb26808b4adcb44e7fe8ecc1bf70a30756ceb5`)
- go@go1.26.10/src/internal/syscall/unix/at_sysnum_linux.go (source SHA-256: `f6ea8f6c906a1d476ce6327d037b186ceb1c8228e80e9bd867a7ea9bc040600e`)
- go@go1.26.10/src/internal/syscall/unix/at_sysnum_newfstatat_linux.go (source SHA-256: `1c15f30d713ca93199e4cf616ac3582a23d96d1e55d44902a39e0c43655a1c0d`)
- go@go1.26.10/src/internal/syscall/unix/nonblocking_unix.go (source SHA-256: `3871841b05dbc99b8d0d6cb51765e2deea09ecd1c394b95e67a52de068db80c2`)
- go@go1.26.10/src/internal/syscall/windows/symlink_windows.go (source SHA-256: `2859cbcb9f1c2aeed58de4427d60d4fbc7626fd4c73bd261d84e997cb7317cd3`)
- go@go1.26.10/src/math/trig_reduce.go (source SHA-256: `fa54b4e9e792707630cd402f19e742ba23c574dc6255867a395e737494a61595`)
- go@go1.26.10/src/net/error_unix.go (source SHA-256: `bbcd49442466a102c8b1d0963bed12afb6cb866cbbec8144f66ded21631631c4`)
- go@go1.26.10/src/net/error_windows.go (source SHA-256: `55c05a49e3bc8a0b7ca9c806848e3fb3847ec07c3864ffafc1a92fead6bac0a3`)
- go@go1.26.10/src/net/http/roundtrip.go (source SHA-256: `a02501dcc615b6af4144655cd28f6a6901973267d2a188794e4552104dcf589a`)
- go@go1.26.10/src/net/sockaddr_posix.go (source SHA-256: `59cc9d08c7a6f916e438525a3f9ee2ade69ded11e4ebf37d0d05e48eb7f0c985`)
- go@go1.26.10/src/net/splice_linux.go (source SHA-256: `a7fa316df36ea565afea6026cf826b023d458733fea433077e5e59c17323a48c`)
- go@go1.26.10/src/net/splice_stub.go (source SHA-256: `06154ddb6b6ea1f4dc0e07bc10f43a86162619e796208deadff9d26884e8e0ea`)
- go@go1.26.10/src/os/rawconn.go (source SHA-256: `cddcedb3b517a8b1fd1a8e190d24bf11978986723a69c8f1e4a231d03e35a07f`)
- go@go1.26.10/src/os/removeall_at.go (source SHA-256: `77d5ea8558ccf4d1018ac6baad111d67d6bc39ab120cf6efe18065700b6d14eb`)
- go@go1.26.10/src/runtime/cpuflags.go (source SHA-256: `3879433296a7402ae37a718371afee50e016757dd2b5bee22f562ddb5683a94e`)
- go@go1.26.10/src/runtime/debugcall.go (source SHA-256: `aca204d0986b11cc4bf04fa0b402ebfb788701195d581fdc3597367fcea40589`)
- go@go1.26.10/src/runtime/map_fast32.go (source SHA-256: `fdce83b7862d1a746f524b5892d4c56312cc534ac9620352bacc428b1a3f714a`)
- go@go1.26.10/src/runtime/map_fast64.go (source SHA-256: `60f77c9c503734b2ec0c9418c2b0ce9f1ad42fd5fd108659834f390dd994ef45`)
- go@go1.26.10/src/runtime/map_faststr.go (source SHA-256: `d7e8512385e72ed212d37422a4a4fffc85eb979c4d45b56f1cd1e0dc006128b3`)
- go@go1.26.10/src/runtime/mem_darwin.go (source SHA-256: `5561c1a69491bc5c08075764db75364ada35416157fb3650c25aeab00d07aec0`)
- go@go1.26.10/src/runtime/mgcstack.go (source SHA-256: `c25628e388174f391c889daa770731336f6240900fe4c668af7ee0608cd8c090`)
- go@go1.26.10/src/runtime/os_nonopenbsd.go (source SHA-256: `ee262ac265c43fa6826127530808e2b316bfb86bf3d226dbbebea428266f2b8c`)
- go@go1.26.10/src/runtime/stubs3.go (source SHA-256: `8cbbe3b1cff1e5265bfe4f54184a665bc12ded904f5d923a4aa8beeaafc681ca`)
- go@go1.26.10/src/runtime/stubs_amd64.go (source SHA-256: `1c7a7c95d3ba4575b61e47e72bce63b34aafcd7482d974ed1b7313428cf3649e`)
- go@go1.26.10/src/runtime/sys_darwin.go (source SHA-256: `66b83a5ed26a34f31a56f5eed686a146aaf5f4817b5d7a308c5332272a1c2d92`)
- go@go1.26.10/src/runtime/sys_libc.go (source SHA-256: `3504f0ee3ff07263f88859d22004a703ade724f7091566eb56c519d284de84dc`)
- go@go1.26.10/src/runtime/vdso_elf64.go (source SHA-256: `759cda62ff78474fe0ef376c713171e23304bb961c013f0f78b44344a965e0ff`)
- go@go1.26.10/src/runtime/vdso_in_none.go (source SHA-256: `21639cf149c73c080136387c1b26ee886b04997618223fb58ded7ee87ba84938`)
- go@go1.26.10/src/syscall/bpf_bsd.go (source SHA-256: `3c8b2c0e8d280b06a83be945f52c4731d34928ca6ea37747e776ef125cbdf125`)
- go@go1.26.10/src/syscall/flock_bsd.go (source SHA-256: `10137eb2e36d95d02f3376c43f24ec0a8d1613b5c2e512a29f7faaf1a9aaa840`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20/chacha_arm64.go (source SHA-256: `50974d7653c355d8356aaa33318f1df5d9127b5dc96a38b5d88e91aeb188feb4`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20/chacha_noasm.go (source SHA-256: `bbdb67ceb30ef13efc54e17935f47e00744b96754e03364e606d4f1f8f9085c1`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20poly1305/xchacha20poly1305.go (source SHA-256: `7325ab1f70b6ab4a241ca0237f050202ddb35225a690c8146006e51a6ceeb685`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/internal/alias/alias.go (source SHA-256: `7bb0ef254a7e9360ec43e3129fad11d27edcc611bf87c12dba089f6e381f462e`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/internal/poly1305/mac_noasm.go (source SHA-256: `f5308cd6f14bab1b00963eeae8137d63fc5452db7f7d49276a56f457bfd89d2d`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/internal/poly1305/sum_generic.go (source SHA-256: `b0094a2895d5bda42dcaaf57c0b31fc914c3b6c6aa6237aab0100a6d78346933`)
- go@go1.26.10/src/vendor/golang.org/x/net/http/httpguts/guts.go (source SHA-256: `3da54270a625540809f2f767e54e18acf27199aca64b6cf16d6d3fa3999de959`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu.go (source SHA-256: `013844c4888fea8b3fba1ab6c80e97b0c504df8bb6645c1c17178a9df157ac93`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu_gc_x86.go (source SHA-256: `2eec0a58e170f14b054a9d0c006fede8cd5c075c5d7cde1fbab688b4b793588e`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu_x86.go (source SHA-256: `b3f1841e3f4a9efdd6e2fdd2bfd571a8dae263c6ba34f45fbfaa803cba1afdc1`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/gcm_arm64.s (source SHA-256: `b5f403162a8d298d425b0885f96458d56c44d08999b0afc82b4aa59109589d83`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p256_asm_arm64.s (source SHA-256: `4f4e28c1651f9f0f86b74861ae1843913bb884ddcee3f3ba99714f7c38ea2838`)
- go@go1.26.10/src/crypto/internal/fips140/subtle/xor_amd64.s (source SHA-256: `6a7b554b18e085c2f1dedc8691974322b557eb30ccf3ec4d60f838485d75b854`)
- go@go1.26.10/src/internal/bytealg/compare_amd64.s (source SHA-256: `9e9482d340d3ff9a4b2aae6350fe02e80938d2cffea2ddf7f209d9c77ea93c15`)
- go@go1.26.10/src/internal/bytealg/compare_arm64.s (source SHA-256: `9fa3a81ed874f9a55984f0cd1fb95b5826da2fe24989a6721e743aeb0623ae8e`)
- go@go1.26.10/src/internal/bytealg/count_amd64.s (source SHA-256: `c8a0ddfc1cc2c278d3076003a2b4fc7c7786669cdcee9f782bd9ecfaa0bb07cc`)
- go@go1.26.10/src/internal/bytealg/count_arm64.s (source SHA-256: `8143fb33923aa726dd2ce21c5052b839a75fd632bf9f76066344199e5ec5c0e0`)
- go@go1.26.10/src/internal/bytealg/equal_amd64.s (source SHA-256: `d5c88b17e90ab21582c77774684d90ca244f3224f09d05625de4d9c2d56b2b21`)
- go@go1.26.10/src/internal/bytealg/equal_arm64.s (source SHA-256: `f42527a9b2ee68d0dfe87f52bb0cb69a4bfd67194820e39eb1b8ac4edf2c66d4`)
- go@go1.26.10/src/internal/bytealg/index_amd64.s (source SHA-256: `f475516c44e8a0bcffbef3a08daec4842d42f0f5eb211857e8af0a81d2bdad8e`)
- go@go1.26.10/src/internal/bytealg/index_arm64.s (source SHA-256: `26d775899469f4dc8eed998be6acfb386944e076ef0a50b208533601360f1a9a`)
- go@go1.26.10/src/internal/bytealg/indexbyte_amd64.s (source SHA-256: `362335dc6a582b92f6def9f118dbdbb492194862b322b94dcedde6d90717549f`)
- go@go1.26.10/src/internal/bytealg/indexbyte_arm64.s (source SHA-256: `ff65073dfdaf022aee0444d59e5a4980a1dc483d45544cf9227ddf4fac737a5e`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20/chacha_arm64.s (source SHA-256: `73036d54d1961a9e5cb8ce57298bd7803443dbe0de7cae4a9143b62d3a7d84b5`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu_gc_x86.s (source SHA-256: `74ac7fc7ef9c56c3306238cf031ea8ef7c0312a7116cb9fba6d07f0b1382df80`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/md5/md5block_arm64.s (source SHA-256: `7f36e29c7565659b108431aa4ae28fd16df4adf8d534bc84392448039a4c5b77`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// ARM64 version of md5block.go
// derived from crypto/md5/md5block_amd64.s
````

### Original source notice

- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20/xor.go (source SHA-256: `c3adc8555766e36629b4dd3fcd579969d2d2737c65a379db997253a2b7b18072`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found src the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/hash/maphash/maphash.go (source SHA-256: `42a498909b268b4f1eb1b172a4e29f4fbcf383b1592465a0d119e0e98f4377b9`)
- go@go1.26.10/src/internal/bytealg/equal_generic.go (source SHA-256: `d6c5c106505c972fab2400024c14e04aa7eee1f972826e3b56787ee09f169247`)
- go@go1.26.10/src/internal/oserror/errors.go (source SHA-256: `4c6ff8e18d1fe326e8b27350b60cdb55e79757452ef79355cc5b9f5ddfde5920`)
- go@go1.26.10/src/internal/poll/errno_unix.go (source SHA-256: `3f1ecb57a178ae80e05557700fc66b39f839c1205075bda2b83f1562ca222b32`)
- go@go1.26.10/src/internal/poll/errno_windows.go (source SHA-256: `b84cde89a8559a123850534a75595850dfb8d762a566b9892aabc407c20fa78e`)
- go@go1.26.10/src/math/bits/bits_errors.go (source SHA-256: `1341283973609b16798176456337e3a0cbff9bfb32f79673a45197223471d7c6`)
- go@go1.26.10/src/math/fma.go (source SHA-256: `e73845c65d1cfc98669fe58d2af47269727d41b554475056f66c41a4f21d0669`)
- go@go1.26.10/src/net/http/clone.go (source SHA-256: `8068d8fcac423a8c498dfdea4a6825736abd80776565c418139348336fd7e0c8`)
- go@go1.26.10/src/os/error_errno.go (source SHA-256: `3894fbb36a02dbdb543c074a2d7d316be881861dd787bb7eec06a3c074ef650b`)
- go@go1.26.10/src/runtime/cgo/linux.go (source SHA-256: `31545deddd648458011c4b21a9c51eef1f180eafe5d54f8f2866c3fa3c4984f3`)
- go@go1.26.10/src/runtime/checkptr.go (source SHA-256: `f38f5a588ddfcc3d528e7b807e5176adb2abc8a82469110fce783e432879cf7a`)
- go@go1.26.10/src/runtime/debuglog.go (source SHA-256: `6f450a47651e81c6cc175cd7ef521d0dc95bfb5b7b2499a473f583c89231be28`)
- go@go1.26.10/src/runtime/debuglog_off.go (source SHA-256: `2403e8ab6ce05543e5ac6e54120c6f21733aca66e99d5f3be803f20933e9f649`)
- go@go1.26.10/src/runtime/mgcscavenge.go (source SHA-256: `809c6b1738c31d8e6acaee363fbb73653c895d56b2284b89b2dd085f2fa12f27`)
- go@go1.26.10/src/runtime/mpagealloc.go (source SHA-256: `63b48da72491a71b0f4103d2f27183f17732d6d2f9b4c9360dde1e97cb3af36c`)
- go@go1.26.10/src/runtime/mpagealloc_64bit.go (source SHA-256: `7b813530b97d23615798cd2359170ee9e1afd8c6e4ffc93d53903f0db0727a75`)
- go@go1.26.10/src/runtime/mpagecache.go (source SHA-256: `ad1427cbe7bb0a07746323b403093f484f750dfdf6f89df876cbe5720cec9c4e`)
- go@go1.26.10/src/runtime/mpallocbits.go (source SHA-256: `8996a481943b5054af0962241b3f16f8ac13d79b5de7b62af904213a7f8b6ccb`)
- go@go1.26.10/src/runtime/mranges.go (source SHA-256: `a9a4f58b3d28a047b9b097b507833f78c96b085e48ae6c8b35bd28db97cb1aca`)
- go@go1.26.10/src/runtime/nbpipe_pipe.go (source SHA-256: `8045e3f11b9ad656641c985b388f2e7c789fde1b921e0d8d4356ad61bc860447`)
- go@go1.26.10/src/runtime/nbpipe_pipe2.go (source SHA-256: `29976eeaa240c5b1ff32866263a1e701831b6fc70a06bf79d6aabc999ab45e4d`)
- go@go1.26.10/src/runtime/pprof/pprof_rusage.go (source SHA-256: `b6f1218e14f49794e7f903da0b90e5430d3b6ca3c080005bbedbf3cf62360493`)
- go@go1.26.10/src/runtime/preempt.go (source SHA-256: `c50342a043676dab55b58ea9aabc90a04cfb0c0cdd9a70e475239ba9460f5ce3`)
- go@go1.26.10/src/runtime/sigqueue_note.go (source SHA-256: `84a465db62b4f99d10ecf64310a081840d22727bc1c0a5203e5c4dba627d6a77`)
- go@go1.26.10/src/runtime/stubs_arm64.go (source SHA-256: `25566bbb705b1057a940b382bfb8a3a2f728ac05feaa09e544d7826533b50c42`)
- go@go1.26.10/src/runtime/time_nofake.go (source SHA-256: `d42faad51f376b9d43b4944620b4208a24bcb48ce7c0d9ef1d7786e3467b1f42`)
- go@go1.26.10/src/sync/poolqueue.go (source SHA-256: `25b900ffeb803a9ecd4c1926a5faaa3ba8d7b4adcd71bbfbd373cb8606910dae`)
- go@go1.26.10/src/syscall/sockcmsg_unix_other.go (source SHA-256: `a2320d4a54bb9b491374508605c892597c85e5a07c57574a2aec4ed7885ae021`)
- go@go1.26.10/src/syscall/time_nofake.go (source SHA-256: `8c9c3e97fa48debbdb2808cc290918e2abc49fd9198811fd8e69477d50467ec0`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/byteorder.go (source SHA-256: `825146fd4557b1cbd8161fa28bb4be8820089848d695316edeecb7fd5a551f8a`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu_linux_noinit.go (source SHA-256: `620abd97199a2230f59165a2bc58c3ba780b6801ec0246cfb7319553dd795226`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/hwcap_linux.go (source SHA-256: `4101df793fddf76dfae477f917928008bc4a797446cd4ad44cde6e906c3f8713`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/reflectlite/asm.s (source SHA-256: `f92c821a60f978ed528d1bfe5846212b5047eba1628f7735abcd95e9f444e9e5`)
- go@go1.26.10/src/runtime/cgo/linux_syscall.c (source SHA-256: `fb53e53bc37d5ba8f8df641e4f6e0abd705e3eae7815ccc87f55911d20c36133`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/ecdsa/ecdsa_noasm.go (source SHA-256: `066437a7409502d153efb49449f77114b8b84960227c474bf3d01358807875db`)
- go@go1.26.10/src/crypto/x509/internal/macos/corefoundation.go (source SHA-256: `a9170ac1609926cfdd2203e96b4b61c4f4cc81d93f322284e93e9b5e39c40702`)
- go@go1.26.10/src/crypto/x509/internal/macos/security.go (source SHA-256: `1ba9ecfaa2764ecf63f324c8fc424b06a7191e61da1f94e30f62c78af8c4422a`)
- go@go1.26.10/src/crypto/x509/root_darwin.go (source SHA-256: `b6d9c56fffd642198a684b14c814b44a7c8c4a63d13db4d7606a8d62aee66c37`)
- go@go1.26.10/src/embed/embed.go (source SHA-256: `76a63eaa0fb982c485df5bd3cc59e4765397f5f1103d2420ce45ddbf8380c52a`)
- go@go1.26.10/src/internal/abi/abi.go (source SHA-256: `0da16138388dd2675ca8aac385667fde079780fbce86c390a3c24cf7e7663493`)
- go@go1.26.10/src/internal/abi/abi_amd64.go (source SHA-256: `b8d763278ba88322c8a7eb1b0aec4ec04063e0a86c2940eeee865a04b12b78c7`)
- go@go1.26.10/src/internal/cpu/cpu_arm64_darwin.go (source SHA-256: `4c77155ef5a80e71603f749fb4c2433f381ed9f706ffc0380c00f337d5e73e24`)
- go@go1.26.10/src/internal/cpu/cpu_darwin.go (source SHA-256: `5c9dafb5c15b5cb405520aa711230231111de33601eb732922991fe008ae2384`)
- go@go1.26.10/src/internal/cpu/cpu_no_name.go (source SHA-256: `8215f72c7d9df7b7b3b53f1d6a5df311816e57cf6e5a06f21f7be50097ca79cc`)
- go@go1.26.10/src/internal/poll/copy_file_range_linux.go (source SHA-256: `39829503f0b04fdb7ccbfa7928f0faf87bffa6143ea66a11ecd2352d92e7cb04`)
- go@go1.26.10/src/internal/poll/iovec_unix.go (source SHA-256: `17ce8fe36c9a5b77aeb11559bb8b22cfdb15d2626fbe1c7d2e50b6ecf490c464`)
- go@go1.26.10/src/internal/runtime/atomic/unaligned.go (source SHA-256: `4c788293b984be07a34492f389d27b707b97016e8e5c6eef2b055f0e800d05e0`)
- go@go1.26.10/src/internal/strconv/atoc.go (source SHA-256: `b64c4e2aab7fc61aca026d5709a8be703cde3dfd42c78f25489e4687e41785bb`)
- go@go1.26.10/src/internal/strconv/atofeisel.go (source SHA-256: `af5b2e23d8ab9624febb8751364b3c48388e02712817db527dd4abfb89ed7bce`)
- go@go1.26.10/src/internal/strconv/ctoa.go (source SHA-256: `12e247aedfb8b2ede713f28aed0717770463e103886d63d595bf8fe09f445155`)
- go@go1.26.10/src/internal/syscall/execenv/execenv_default.go (source SHA-256: `42dcdfc64c4b31806fb848c61f3c6efceff7659db131964240ccf57813355af3`)
- go@go1.26.10/src/internal/syscall/execenv/execenv_windows.go (source SHA-256: `e636f06a1d6877c65a84ec60ae76416e5095ea9b2ee244fc48cf7b75f87156f2`)
- go@go1.26.10/src/internal/syscall/unix/copy_file_range_unix.go (source SHA-256: `7fbf814a8a7d4fbd0c76d1827e50e8864b4bb701cad01277b96af9d0826809e0`)
- go@go1.26.10/src/internal/testlog/exit.go (source SHA-256: `8f82410c509c4e38d801c9a022368050e878368400bc6a678b62e62ca5f1215c`)
- go@go1.26.10/src/internal/unsafeheader/unsafeheader.go (source SHA-256: `aeb2bf9d6c18b361ef4d2830dd63675960268432a4c5dd6131a374d654094da4`)
- go@go1.26.10/src/io/fs/fs.go (source SHA-256: `8dbcead59f3be21f03691b92b77b8d4430cefb63dd39736f773c492e9265f83d`)
- go@go1.26.10/src/io/fs/glob.go (source SHA-256: `c3d08c0f3207f4bcc1703127c324083c9ac65a0e2fa685379f3df82c858ab399`)
- go@go1.26.10/src/io/fs/readdir.go (source SHA-256: `9709461e50d64cfd0a670382912aaca59af04f0506df809b2b56ab50c438242e`)
- go@go1.26.10/src/io/fs/readfile.go (source SHA-256: `69814861773bc26f58d1fd49167b30f39d148475d62e279c78a290c0b499c5ba`)
- go@go1.26.10/src/io/fs/stat.go (source SHA-256: `f2e004297868594204e85f903103f64c05b1e729066eb6b97740b0fd642fdc12`)
- go@go1.26.10/src/io/fs/sub.go (source SHA-256: `b067c9b86e64e3a0b282072baa95d0e600085584c207a116fdec7803b15237cf`)
- go@go1.26.10/src/io/fs/walk.go (source SHA-256: `09de459e421e28f05c2eae1b5848cc616aa0879760d24e955faa056de47c26f0`)
- go@go1.26.10/src/net/fd_posix.go (source SHA-256: `31e1ad8e0ed8b6a86d9f547c9695a1d67b35a553e0e251ea66362fb22c5354f2`)
- go@go1.26.10/src/net/netip/netip.go (source SHA-256: `dbebdcbfd55f9619d57f55910b2af6240adb84dcd58bc76b9b396e479b91815e`)
- go@go1.26.10/src/net/netip/uint128.go (source SHA-256: `34cb2ed4ee74833234158b4f1e3a052713cec5657af1fa18e2284e70c6ee5604`)
- go@go1.26.10/src/os/dirent_linux.go (source SHA-256: `554ab5920f7606b10cd63c5be7146fce21255c56e52bb44dd994c943d0962f71`)
- go@go1.26.10/src/os/zero_copy_linux.go (source SHA-256: `6cada0ec33a15bc5be2710cf5404977d774525dc5150071d08f1e05f104c5e61`)
- go@go1.26.10/src/os/zero_copy_stub.go (source SHA-256: `5ff33597850b18e76f6fb6b004b55e5430c4316ee627261330368afb22049f42`)
- go@go1.26.10/src/runtime/cpuflags_arm64.go (source SHA-256: `d13970d6538f225ee81146f7bf019a4e492f5bde1e977d7ff95be8c449c0ad51`)
- go@go1.26.10/src/runtime/histogram.go (source SHA-256: `7e9c532f80d134298db85fee79a1e251de00507d996976f2ff62060c45486818`)
- go@go1.26.10/src/runtime/lockrank_off.go (source SHA-256: `7939436cf311332dfaf11dad3d4ad42176e549409bd735f1c6eb2e0e264cea5a`)
- go@go1.26.10/src/runtime/mcheckmark.go (source SHA-256: `fa147698389ed989047a0f10f64823117e55a023d2df5764c385cd8422583672`)
- go@go1.26.10/src/runtime/metrics.go (source SHA-256: `880e299fcccadabb5ae36ba13201f4efd02bf58231d82edefced6b5a58d50cce`)
- go@go1.26.10/src/runtime/mspanset.go (source SHA-256: `85faad45efe16dc88c72253e6279f58e2cb161d948534dda5c47672243ad367a`)
- go@go1.26.10/src/runtime/preempt_nonwindows.go (source SHA-256: `4aae15044bdee656cbb361b9c06c9ab2f14b25d28c4ba6c407e27c225e7a6d2a`)
- go@go1.26.10/src/runtime/sys_darwin_arm64.go (source SHA-256: `35411e5f1d4f9ac0bbe73603a939970aa1fcaff849c90d0a9b608163ca7a3090`)
- go@go1.26.10/src/sync/runtime2.go (source SHA-256: `5bc4014718a892af6591e26e6c87105a001b7d1d1f93c3ec4c289128c80f6ba4`)

````text
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/subtle/xor_arm64.s (source SHA-256: `487b46b18108844057ab4627a67a4502936afb0091899bccc617871f5acbc540`)
- go@go1.26.10/src/crypto/x509/internal/macos/corefoundation.s (source SHA-256: `7c4c31f01ff02961d2054b21e8bb3581889f009533d37c71cf2448e76fb71c50`)
- go@go1.26.10/src/crypto/x509/internal/macos/security.s (source SHA-256: `1bcdb3e38afa2cbb67357e402019d3302ef100b1da2e444ad5c996a5e3dabbfa`)
- go@go1.26.10/src/internal/cpu/cpu.s (source SHA-256: `39594514f3cb069cecfc6782f81385cbb3570a2ca1b0f3028db1dfa75fd5837e`)
- go@go1.26.10/src/internal/cpu/cpu_arm64.s (source SHA-256: `03dbaac3d31973cdb520257e3b5b8ed33d9afdc0e7ab0c616862557b24300a4a`)
- go@go1.26.10/src/runtime/cgo/libcgo_windows.h (source SHA-256: `d988ed23077a36acf06e83d98e078c176f508b9f751dec437ecc862e8cee8e3f`)

````text
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/elliptic/params.go (source SHA-256: `73064e7cf9ea3a268bc601f77f5bb8fb5550dfbe06eebf52f8ed8643c44d91e4`)
- go@go1.26.10/src/crypto/internal/fips140/bigmod/nat.go (source SHA-256: `02e9175e5315b183eb60afa8d5ebb4c52d2b0d88289b81e16823a8ebb695cd49`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p224.go (source SHA-256: `728051c5690f3442e87d4b9c7e8d3d879e2029f28c291049fbe4747bcaf235ec`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p224_invert.go (source SHA-256: `3d13e6db8c00173af0fc1409d699d227bf33a16a6191f4594c1e3a5a9e319135`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p256.go (source SHA-256: `b5e874ddd907466916c93ed1bca988d7264ee2cadcc4892bcf72c8d3ea50fa68`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p256_invert.go (source SHA-256: `2677f82201b5789d9257a9cc5b5ddd6036858245574d9ee220bdd1715d2e67bc`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p384.go (source SHA-256: `0f8f338749206814b9cd551bc571ae2a69c7d300d66976ddd59f64df15e1243f`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p384_invert.go (source SHA-256: `35f64d7c76a9ed425dcd9925696ce4c69c064d5d1ab9093c3cebca3915e508e4`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p521.go (source SHA-256: `05a2e3778e9601f8d3e10df3c367eb993bfe3a975e7258751d13321ca1512a24`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/p521_invert.go (source SHA-256: `f1cb1bd68ebcb1e147e493298f5a33a06c22ae8d278dacb84ddade9680b5c8ce`)
- go@go1.26.10/src/crypto/x509/parser.go (source SHA-256: `70ae7c65f68d17a59e7bacb0d1cc520e8e96e8eafee2b7224465180437ab3a35`)
- go@go1.26.10/src/internal/abi/abi_arm64.go (source SHA-256: `f9a322657c339f34ac6be9e12f967922aaa1fbc19b1c159f330ab2af769c5b77`)
- go@go1.26.10/src/internal/goarch/goarch.go (source SHA-256: `4ba164d62956cbf38638637453837789c39761e739aeff267702b33bdf5908a9`)
- go@go1.26.10/src/internal/godebug/godebug.go (source SHA-256: `cd492e76ae63872ac4310fc09b2e097229588dbf114d14d713be3f027b0233c8`)
- go@go1.26.10/src/internal/goexperiment/flags.go (source SHA-256: `b34bfa58c283c17def13d43173fabd633dd361ebc22a1fe172750ceb14c43a3a`)
- go@go1.26.10/src/internal/runtime/atomic/doc.go (source SHA-256: `9e0cfe8821fb34bd0511f55b7e269c6b74dc9956e5d2b2db367349d1aa7695b6`)
- go@go1.26.10/src/internal/runtime/atomic/types.go (source SHA-256: `5296242a02223a60f714e9358ea58aae52dbef7a4fdf1f90cd479b86360d0c70`)
- go@go1.26.10/src/internal/runtime/atomic/types_64bit.go (source SHA-256: `1ea374803ac8f483de8a3efa9e065157a7d2352652a4ac7bcb0dc17756daffb6`)
- go@go1.26.10/src/internal/syscall/unix/getrandom.go (source SHA-256: `1cafffb178e0a2182349ef5a39fd89aa01f2a9d1fd565afdb947416b9ab13e40`)
- go@go1.26.10/src/internal/syscall/unix/net.go (source SHA-256: `866c1318e10b6d9d14bddd6a403220a69e69b52549510b63b3365d629e75b3fb`)
- go@go1.26.10/src/internal/syscall/windows/net_windows.go (source SHA-256: `2dd1e1e69816acd2cfcc5a679082504e5e96a3363024c824a3c32d71de3ec998`)
- go@go1.26.10/src/maps/maps.go (source SHA-256: `0f69acac135baa8b644e2479e514835b4d991bcc894916c2e34b8739a171fce2`)
- go@go1.26.10/src/math/dim_asm.go (source SHA-256: `15b346a46f5926bb90b54de6c6c3de73c55b4ad462ab69ed168f23ae1c1470f0`)
- go@go1.26.10/src/math/exp2_asm.go (source SHA-256: `fdc3d87224285e8b57e8bf7c61cc0881fbf735ffc4df54841103a5b314cc7cf4`)
- go@go1.26.10/src/math/exp2_noasm.go (source SHA-256: `199345cbae2f94fed9a9a2160d8e1a15389daa958ae92621c6adab22208e114a`)
- go@go1.26.10/src/math/exp_asm.go (source SHA-256: `152874757ad993e65449015bcc5e7a494c8b5139e6774f89b9c6c43a15414cda`)
- go@go1.26.10/src/math/floor_asm.go (source SHA-256: `8f537d2f130c98357554f7b7aca14c5dfa89614e76c8a10cfdf0bdfba8056afd`)
- go@go1.26.10/src/math/hypot_asm.go (source SHA-256: `bb10755f29d0625ee52cbe056320a5fbc9a819f5b6e6f9462259cc9cd28bb96e`)
- go@go1.26.10/src/math/hypot_noasm.go (source SHA-256: `ea186403dd7a55eb45513049920cc7663ebffe80588d43aa6081b2602005531a`)
- go@go1.26.10/src/math/log_asm.go (source SHA-256: `92751cae986507de5d497b53d6d5459e84b4afc22e83cdd504e13d99997dfb74`)
- go@go1.26.10/src/math/log_stub.go (source SHA-256: `c945334a77924e23b19533d80e156e0286bbf55fe4e189f33a0866f9440f689b`)
- go@go1.26.10/src/math/stubs.go (source SHA-256: `191ed067c7368ec75c25d1f7298720e14bbe7280230aa4d6ae5f791a77452015`)
- go@go1.26.10/src/net/http/internal/ascii/print.go (source SHA-256: `f2af00c0c1b74c279c0bd7d1b0054a5440cb4046cbd60292dfc1aa4ada76845e`)
- go@go1.26.10/src/net/http/transport_default_other.go (source SHA-256: `c909fea6b4308f7c164352032a694f18f87105649005a2631d29e7b5d6b5c587`)
- go@go1.26.10/src/net/unixsock_readmsg_cloexec.go (source SHA-256: `c9ba39caec1254b9da4ba2a2559ffb9be5008389f2ece7e1934cf643da35f9c8`)
- go@go1.26.10/src/net/unixsock_readmsg_cmsg_cloexec.go (source SHA-256: `b3767e4a73f1f23365d286d5cb2c0bb38b9a7bae3e14b5a2d1499d7385be8c2c`)
- go@go1.26.10/src/net/unixsock_readmsg_other.go (source SHA-256: `863347fc0f8b78feca2ef31b87894e4c121738121a78a411ac132be9eff20455`)
- go@go1.26.10/src/reflect/abi.go (source SHA-256: `a62034c1299fd4ba11a40ca398345f317038c40bb6c774daa7594ccf12816d9f`)
- go@go1.26.10/src/reflect/float32reg_generic.go (source SHA-256: `225b092a447b002cf4b818e5c01d28dbd1e4d1bb30b681a840aa7f3f7c7781d0`)
- go@go1.26.10/src/reflect/visiblefields.go (source SHA-256: `287edc70dd44cb440aa267eb9ca3164c89c80b25af69abc1589728fbf82ad810`)
- go@go1.26.10/src/runtime/asan0.go (source SHA-256: `34ffa01d044c5abbf40f0fcfa7634e94a7d7320070822fdd874c8c288e114f7e`)
- go@go1.26.10/src/runtime/cgo/handle.go (source SHA-256: `4f30cdd7bd95b2fa6bcbd0c8cf47395b3325e0be642468d4d09aca0f0b7481db`)
- go@go1.26.10/src/runtime/mgcpacer.go (source SHA-256: `b9399df0e9db1b8fb8bc701accfbebfe4615303b1f649ec1c716024847969e73`)
- go@go1.26.10/src/runtime/tls_stub.go (source SHA-256: `1596177a6668e35fea9dadebfa93dbe057a601886081665107e32fcec5fab2cb`)
- go@go1.26.10/src/runtime/tls_windows_amd64.go (source SHA-256: `2e1d6937c5be4c3932fc7d7af4778e2438b7d6350d7cc5da40589ff25bae06e0`)
- go@go1.26.10/src/slices/slices.go (source SHA-256: `8f2f5990f2189d57e530ddb8bf8224bab3e28e12d5da9f44e6d082c94228adcf`)
- go@go1.26.10/src/strings/clone.go (source SHA-256: `5c18eb3b45b7d9ae63d0cc43f9e337c077de05dd47fd5a1ea49dafeb7e1b862b`)
- go@go1.26.10/src/vendor/golang.org/x/net/idna/go118.go (source SHA-256: `70f0410bd261f4c487d299926c0a7cc9e5bf3608fe8d727d8e930fa41fdca89a`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/abi/abi_test.s (source SHA-256: `28a2adea62ff2469766d52f89bfe67674e388a9ecd7c41f04bf4be063517b8d8`)
- go@go1.26.10/src/internal/syscall/unix/asm_darwin.s (source SHA-256: `055e996ea98e5c629ff48e399c3cfe6c18ae22e36b89123259a20a1864b3b0a1`)
- go@go1.26.10/src/runtime/asm_amd64.h (source SHA-256: `85fc2d5bb80e6db36cef4da5e5c02278d0dcaa6103e5660629b4b713cff15d54`)
- go@go1.26.10/src/runtime/cgo/abi_amd64.h (source SHA-256: `7832f3b0a3d91852a53705f3df6d16d69e1e785cfd1d733ee33229411cd26717`)
- go@go1.26.10/src/runtime/cgo/abi_arm64.h (source SHA-256: `ea02dbdbe31f446716c290bbeb5b98f940ab71c1756c5a7da68711ef73e2b3ec`)
- go@go1.26.10/src/runtime/time_linux_amd64.s (source SHA-256: `328970a6fb6973b917bc914931bb32d001f698b6a399b4d8d5f7c1939051f612`)
- go@go1.26.10/src/runtime/time_windows.h (source SHA-256: `f50c5452a69fda66f7fbf895aa73c5d50de0cd24f7b285637fbaa6837c8f300c`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/ecdh/ecdh.go (source SHA-256: `63e3e45ac9791f633e28bebeed63cee570a1868f2aa1064741e8a302ffc79ec8`)
- go@go1.26.10/src/crypto/ecdh/nist.go (source SHA-256: `1c4f8e4e6d809ab78d22a2255b237cfddb2e957b57c07927ef7fb414ff622953`)
- go@go1.26.10/src/crypto/ecdh/x25519.go (source SHA-256: `6dff1337f8540910d0bed609b7acf338e365c26be56df9d6cd60fd6b49afa41e`)
- go@go1.26.10/src/crypto/ecdsa/ecdsa_legacy.go (source SHA-256: `b35d1db7a48be3aa89d1823f0f4d87c938cd850bfef784ebb9c6c3bfe6604790`)
- go@go1.26.10/src/crypto/ecdsa/notboring.go (source SHA-256: `a723fd5b176febb6a5435861b5273ed9cf173f4f18aedabc0f1937806e57efd7`)
- go@go1.26.10/src/crypto/internal/boring/bbig/big.go (source SHA-256: `a0c5d18d5b08f678eeaa05029740acfa2399bf665683e65f57859bbbbc210646`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/nistec.go (source SHA-256: `0ba361874ac7c4c91f95f5be67dd72dd11bab4fa66aebf6e84fd717a051c1ab1`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p224.go (source SHA-256: `e092447bf67f25f347fb0727ce2cf951678799c2685a392eea6df176222efe54`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p224_sqrt.go (source SHA-256: `406c33177e0db01283f57ed849c89e1765ce3b108cc13b1a415f0265ed387376`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p256_ordinv.go (source SHA-256: `241942f92e4d0e3e745f77ea51c642f8687835a6e3e9058418dad26180aced0e`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p384.go (source SHA-256: `ea1fa57656170d1b57d954dc577aa69e40ba6bc1cea8c2e4e13c37ee1af477b4`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p521.go (source SHA-256: `faedbfb701558b518ef0e6c38ec81edae331fb61ac3c8610b64cca2bcfe047ee`)
- go@go1.26.10/src/crypto/internal/fips140/sha512/sha512block_arm64.go (source SHA-256: `cac486830641b98c295fd20d07606b8eadb7d6a2f258d891dbd80f0ebabdfe48`)
- go@go1.26.10/src/crypto/internal/fips140/subtle/xor.go (source SHA-256: `7ee3aa840d46884e788f1fa42fc9d276dd8272fc44f9df3f57bf6467291f13c0`)
- go@go1.26.10/src/crypto/rsa/notboring.go (source SHA-256: `d6b371eb52957ef1adb969e70d7a28cbcb3360837bfd401e6805655ee4d8e266`)
- go@go1.26.10/src/crypto/subtle/xor.go (source SHA-256: `f438a661cf333ea488b3c4b227bdba8d2fd5914d916c7b327eb39e343db27296`)
- go@go1.26.10/src/crypto/tls/cache.go (source SHA-256: `553ad7a40a9b198bea838833e3e5ed5660ec9b031ad75f58a9ea119e44bd3bfb`)
- go@go1.26.10/src/errors/join.go (source SHA-256: `87da7c120f71d17e37f110ae1dc6af6493ec841a22ee825e13a6a8024a014dba`)
- go@go1.26.10/src/internal/coverage/rtcov/rtcov.go (source SHA-256: `896d5a0053f85c64b2dac75934580e39576ba558b94eefe6fba1f66599c4e1a8`)
- go@go1.26.10/src/internal/goos/nonunix.go (source SHA-256: `df24b5ff13513e2d4cf4125b84792c413767951cf05ed55995f5d3e174b4ed68`)
- go@go1.26.10/src/internal/goos/unix.go (source SHA-256: `e1f33ee1746ab75f8fa813645d2e36a0b1dba2bb0b9c5974e7fe658f210b32bd`)
- go@go1.26.10/src/internal/runtime/sys/consts_norace.go (source SHA-256: `f813bc3694d70c2f77b868cb52789fec13efe5c1308f71674ef16edd54e9c058`)
- go@go1.26.10/src/internal/runtime/syscall/linux/defs_linux_amd64.go (source SHA-256: `814cb91b8ede5df93740a4da08d74e52ace283f60d811a3e5dc64df949852216`)
- go@go1.26.10/src/internal/runtime/syscall/linux/syscall_linux.go (source SHA-256: `3efac6fc31e689e3c3073e82aa7f66014fd224ae237e205140eca924c749ec23`)
- go@go1.26.10/src/internal/saferio/io.go (source SHA-256: `efadc03e4126b93f9924cd2d3f70019387bd2945a86657bc6663207d0c8cebc1`)
- go@go1.26.10/src/internal/syscall/unix/constants.go (source SHA-256: `4ccef3971fa8cf67f5ac98c696cf252cf14c2fa56ad18aa5cf4b8998efe03b65`)
- go@go1.26.10/src/internal/syscall/unix/kernel_version_linux.go (source SHA-256: `7e1df29baff4556808f0b4c250456e6e701ec63d2fad5eee5042f1b3aa10c7ae`)
- go@go1.26.10/src/internal/syscall/unix/kernel_version_other.go (source SHA-256: `f05524625bdffc319d4cb41d92d16c9a00f499d602a91d3e4275b3dd0e33cb2a`)
- go@go1.26.10/src/internal/syscall/unix/net_darwin.go (source SHA-256: `1138a6d5f76d72b9c9812bfa86a4b503c1129967e82fc28e33af282b5564f85b`)
- go@go1.26.10/src/internal/syscall/unix/pty_darwin.go (source SHA-256: `51920b16a870dd6fe60be1e2a1db99295dfc04892dc0f3d2fae4af65f14d887a`)
- go@go1.26.10/src/internal/syscall/unix/user_darwin.go (source SHA-256: `9fa7329516fe447d9e42fbd4ae0473ed18d5aec6091e0fa7489fa8a13dd9626a`)
- go@go1.26.10/src/log/slog/attr.go (source SHA-256: `963ad129a6e075c6ffa9b1fa92ce22ab29bed6262f7d5d2f21d0be1aa0fcb3e0`)
- go@go1.26.10/src/log/slog/doc.go (source SHA-256: `4f23bfd50361d93558325ce4d81b4c38b265d64a22329df331986f788bf505d2`)
- go@go1.26.10/src/log/slog/handler.go (source SHA-256: `c8d5520336772523a79b8044708274d787884ff33ff023dc199f08b9d2ebc6c8`)
- go@go1.26.10/src/log/slog/internal/buffer/buffer.go (source SHA-256: `d9a1e570503fafafb16d95429f22b657d5eff456ef48d40bdf5724c5dab5fead`)
- go@go1.26.10/src/log/slog/json_handler.go (source SHA-256: `a10503560615180e018a44a0ebd5ab0c4d0a79f74ab7ee0248ca24277da259df`)
- go@go1.26.10/src/log/slog/level.go (source SHA-256: `ee656ae84acfab6e0bee2a92ef42cc4c8869acbd9fc57a5e6ca7e98071e63a21`)
- go@go1.26.10/src/log/slog/logger.go (source SHA-256: `7789f4aa1da8c627dc84a25d7386213f1f6ac7f0c95ea8ee1f588500a52aeac6`)
- go@go1.26.10/src/log/slog/record.go (source SHA-256: `f816587dae9606117774f1888d877abb950c4f24bbfdf8740c1a92a354f46e6f`)
- go@go1.26.10/src/log/slog/text_handler.go (source SHA-256: `e754f5239c660d425fbcdddeae31df88a579b96c6d2413170b09ef895f229ddc`)
- go@go1.26.10/src/log/slog/value.go (source SHA-256: `ef42323f3c07bb645a5cf86ce7211fe47b724e32c0898cea7b94c2dba02e4750`)
- go@go1.26.10/src/net/cgo_darwin.go (source SHA-256: `3bec0d4702a2416809d4e4ef9aabeeb678d78978a8af5c2a860297f6d3abbe3e`)
- go@go1.26.10/src/net/cgo_unix_cgo.go (source SHA-256: `3beda011e6a0c1dcd5b0c1d1c99518dfd9cb9747487118504faae19237d2dc87`)
- go@go1.26.10/src/net/cgo_unix_cgo_res.go (source SHA-256: `5814ebb62a461d257956c789b0284f9b27efe32836557dc1869db32b29bd5ba8`)
- go@go1.26.10/src/net/cgo_unix_syscall.go (source SHA-256: `29ad469af458d8407c8eccbf86edd64ba13f36d0688e773bfd8100647fc35f36`)
- go@go1.26.10/src/net/dnsconfig.go (source SHA-256: `c78109fe3c4ac2d8d4c84bd2825c9b0d236386a8600e6b5a313afdc0e8914da8`)
- go@go1.26.10/src/net/dnsconfig_windows.go (source SHA-256: `0894f9dfbe6d4669340039ee81436b5b24ee6911a7197abdc8b0edcf71d50ff2`)
- go@go1.26.10/src/net/http/h2_error.go (source SHA-256: `7649362439a44d697c7fcde73dba2fcc95056c7a157e1023cec7b4de24ca1a89`)
- go@go1.26.10/src/net/http/responsecontroller.go (source SHA-256: `7276f45a5fed50dad368b75afe916b3e4dc45bca2e23e438f1095b0fdc6f4ad6`)
- go@go1.26.10/src/runtime/arena.go (source SHA-256: `6853cba54c24b766409f3523774a537bc2b562d5083638d3d86d8de4938badd8`)
- go@go1.26.10/src/runtime/covercounter.go (source SHA-256: `4e2ede89b4fe9e9ea8429a773db441a8dab3cd4962159f842080e97d2ef8d64d`)
- go@go1.26.10/src/runtime/covermeta.go (source SHA-256: `765117a0c617fe622ac461bf1ee7dd3ed075d95ae80cdfc6bf91d280a5609300`)
- go@go1.26.10/src/runtime/create_file_nounix.go (source SHA-256: `e96ebe35523a5982a219aa9fba5f3490605abf5b95a398e427bf78cb82efc370`)
- go@go1.26.10/src/runtime/create_file_unix.go (source SHA-256: `96f29670074191bc32702674ea155fc105b2e4175fc48fa1a5a170662f8cf9c5`)
- go@go1.26.10/src/runtime/mem.go (source SHA-256: `1b35e5003576c3c61fe5b4273b2c2314907a0ebbf0109a1078e33ba44cc3762e`)
- go@go1.26.10/src/runtime/mgclimit.go (source SHA-256: `fd175a6304a6d4e1f26822cdaa0f94d4959b2d1f5c315f6a0b1c342bda8ae6a7`)
- go@go1.26.10/src/runtime/os_unix_nonlinux.go (source SHA-256: `310c8635c507c7e6def6237e4eb323b48443acd9fc84d0e6568b735cb8aa1633`)
- go@go1.26.10/src/runtime/pprof/pe.go (source SHA-256: `bf35baf1f818710a4ece483144ef5860dfa1fd1ee3c1b17f4a3e8bffee41a12c`)
- go@go1.26.10/src/runtime/pprof/pprof_windows.go (source SHA-256: `3b8dce2bde03c35ae89dcfc7616a14a9a649b831ee07e2301217cf539b54a95c`)
- go@go1.26.10/src/runtime/pprof/proto_other.go (source SHA-256: `8a291fc42634b53289e1ead9f008c6eabc673e05c7a3261877aaa541c295bfd8`)
- go@go1.26.10/src/runtime/pprof/proto_windows.go (source SHA-256: `ddb46ad069485eaa1d9e82eae5daaf3e0b0f5564a2642335ddb3a5102d4956c0`)
- go@go1.26.10/src/runtime/retry.go (source SHA-256: `64afd0c6d11790a6324caf058795ca4c61295d03b947637165a764b8d34c0b13`)
- go@go1.26.10/src/runtime/stkframe.go (source SHA-256: `29bfc22b96fdaa16c25c0a1db9778a56f6daa3d4f15622dc1d9bc4faced5ed53`)
- go@go1.26.10/src/runtime/unsafe.go (source SHA-256: `129ae3ace83e83b36edabaa2e3e5ac168475cd8d9a3cc02c81a08b1438b9f98b`)
- go@go1.26.10/src/slices/zsortanyfunc.go (source SHA-256: `50398783e085ac9cc1d7c6249ebe75a5cc4425dcdc6a4cc448ea6e7179126492`)
- go@go1.26.10/src/slices/zsortordered.go (source SHA-256: `f6d8d5df558075ca8299a89f9b5beda3ca7d1ddc7858dc2e52eee44e4dfdd7a9`)
- go@go1.26.10/src/sort/zsortfunc.go (source SHA-256: `35cdbcbc3a63372e12eebc815c249adfd420e9d783a99b0405998a7ccb00ac1a`)
- go@go1.26.10/src/sort/zsortinterface.go (source SHA-256: `8978a49cc4174b0c35a84a2ddf73175d41254eb340ebeda29f21d3e546088a3b`)
- go@go1.26.10/src/sync/atomic/type.go (source SHA-256: `448abbf2d0dfe655fc8d096d2aed2fa37e43b10558c2132f1d0b17c8e4de69b9`)
- go@go1.26.10/src/sync/oncefunc.go (source SHA-256: `0853c1ad885c593d2810e601c495f00abe197b10d12b679f94badf3efe5fa365`)
- go@go1.26.10/src/syscall/rlimit.go (source SHA-256: `f0eee111b4a1fd060d21d2a18878353242c9c9440da5aec97d74667ef634eb51`)
- go@go1.26.10/src/syscall/rlimit_darwin.go (source SHA-256: `de12284de84f8bb3d835d26410c57960f8be0c013ebfe8e37715db0dd64bd5dc`)
- go@go1.26.10/src/syscall/rlimit_stub.go (source SHA-256: `0a1d3dadf3b0f1eb85765f5e7ff1c94322d75c14ffb33f43bf2e229b5c7a8f93`)
- go@go1.26.10/src/time/format_rfc3339.go (source SHA-256: `14b2a58fa295cca1a4f99467f5a5059418ed266a2298316124af97ffa014f271`)
- go@go1.26.10/src/time/zoneinfo_goroot.go (source SHA-256: `984267b22b58101becfd3612f836c3e32a6f4190e149e257315a4179ec279992`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/parse.go (source SHA-256: `bb57db6688bac4a712dbd888a4bf15a1ba2b278f11fcc291c9bb590429a11be3`)

````text
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/sha3/sha3_arm64.s (source SHA-256: `a4f85c44459be7570b3d1bfc361cf75d6cf04a6dd5b8ef905ba32a29c5b39525`)
- go@go1.26.10/src/crypto/internal/fips140/sha512/sha512block_arm64.s (source SHA-256: `e8fe022552851cd595b76c7ae76ea9e0d386a7490f7a815813d46837b0d4bb7c`)
- go@go1.26.10/src/internal/runtime/syscall/linux/asm_linux_amd64.s (source SHA-256: `dd6191356bf0c18b3c9862b19e4014f06e217987b225823612b6da56fb6e193a`)

````text
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/cmp/cmp.go (source SHA-256: `9d8e18366a1a20d2a1981b7f6a0fe53362e801d9ca20eeee56d65536124dd91b`)
- go@go1.26.10/src/crypto/internal/fips140/aes/ctr.go (source SHA-256: `3db57d363c1e78e2d796ae31d4d8b51b892c467bf34fc5364764cc9d4a16072c`)
- go@go1.26.10/src/crypto/internal/fips140/bigmod/nat_asm.go (source SHA-256: `2cbf10e43db0371bd3466068580d3b69960a9be8b57d857c8e66b888fec10bf0`)
- go@go1.26.10/src/crypto/internal/fips140/mlkem/mlkem768.go (source SHA-256: `7387d3f2772ee24ef0e1b35991785669f118c4c7e5720afd923c0212c8fdbd1b`)
- go@go1.26.10/src/crypto/mlkem/mlkem.go (source SHA-256: `710c6cffef8e809e16263797808f9c1dce3b746e733ac12d3bec47c5ee9b604e`)
- go@go1.26.10/src/crypto/tls/quic.go (source SHA-256: `2153e8cb82cbf713afa0a946ccdce334e4fd8ae1d5aa2244b75b6f90ba9d9ed1`)
- go@go1.26.10/src/crypto/x509/oid.go (source SHA-256: `9e9e2f002d23bf6df9b5b0145037f390bca47180366696408ab80a2615dea2b5`)
- go@go1.26.10/src/encoding/binary/native_endian_little.go (source SHA-256: `b031ac2a0198ad9091aba55eaea935b258fac69e2025502aa3d50f194c12918f`)
- go@go1.26.10/src/hash/crc32/gen.go (source SHA-256: `b6f5dd452079fe27e1789a41dee81b24b2b933a655fcaa4ca530be8bafb2fe63`)
- go@go1.26.10/src/hash/maphash/maphash_runtime.go (source SHA-256: `f953ea8eae374bda0cce77c08e63826dc5787fa72182c49edd03a59bcba127bc`)
- go@go1.26.10/src/internal/abi/compiletype.go (source SHA-256: `0988077c09f3c8177032b11d734b0ef4f22b22e9efd27901bbce910e39223f10`)
- go@go1.26.10/src/internal/abi/funcpc.go (source SHA-256: `40daff47ea7abc00e83c99c6256ed7c3af618589568d477125c050cd36fb9af7`)
- go@go1.26.10/src/internal/abi/iface.go (source SHA-256: `a4d66e19e604838666eba4c6c71734eb5485d10983d386f9e0eef48edd71ebbf`)
- go@go1.26.10/src/internal/abi/map.go (source SHA-256: `e3b5f05fab86b2002ea4884f998e32b84093e8569bc68d9976d59faa1e80c554`)
- go@go1.26.10/src/internal/abi/stack.go (source SHA-256: `f21cd3ca6d133b7ab87010b42aa03ff0970742e1b3f1593dd33a5fdb63ec4edf`)
- go@go1.26.10/src/internal/abi/switch.go (source SHA-256: `d2b5f220dd726af671f4f04bd5851188d280084a2d55f46df6ab623d136094f6`)
- go@go1.26.10/src/internal/abi/symtab.go (source SHA-256: `2f3d299e58f35669206f58a0e637a27e3d2a49a5e07dbbe8ad62a0155051fe91`)
- go@go1.26.10/src/internal/abi/type.go (source SHA-256: `b7245bd743802e2b5eced71a1f87963c395c2536cd7ea1b77f2364d20a79bf3d`)
- go@go1.26.10/src/internal/bisect/bisect.go (source SHA-256: `978d518986ba6867f5f96356531f35bfa59b865fbb363b98eb9052e862a2160c`)
- go@go1.26.10/src/internal/bytealg/lastindexbyte_generic.go (source SHA-256: `632273eb2537352d9c14c6f37bf531bf95ff98431e212a24ce3777f436331c2d`)
- go@go1.26.10/src/internal/chacha8rand/chacha8.go (source SHA-256: `6c460d84955fac536cbe55935c19d6c3df77b25122018013fab90c7f78731698`)
- go@go1.26.10/src/internal/chacha8rand/chacha8_generic.go (source SHA-256: `3202dcb9bb0750e7ec3ecff529a09331a9145a090a8c3b5fed57ff5c34e602b2`)
- go@go1.26.10/src/internal/filepathlite/path_nonwindows.go (source SHA-256: `11ed265feded7213c38bff7b28ed9924fd4d664f81ed5cb408f76d5aed9d8fc5`)
- go@go1.26.10/src/internal/godebugs/table.go (source SHA-256: `c8d5b6e7358fd39befe8df0cb31bb859804ced8182b85f0b01355e3a845f7fa9`)
- go@go1.26.10/src/internal/poll/fd_unixjs.go (source SHA-256: `ed73269460d0c3dcd551a2b9caa5c58e5c818a064d2971f8d63b3b13c6fdebe0`)
- go@go1.26.10/src/internal/syscall/unix/fcntl_unix.go (source SHA-256: `5ab1dd517d3ed5beb706bcff14ca89acb6720c1e92a13088e54bffcedc715a6d`)
- go@go1.26.10/src/internal/syscall/unix/pidfd_linux.go (source SHA-256: `c83da30c5aed4890792d8612c78f430c7ab6f3311e8fd73792c24e6636937b4c`)
- go@go1.26.10/src/internal/syscall/unix/siginfo_linux.go (source SHA-256: `e774621eefa370361d1245adf058d8a1752b0bb78ecdf8939d535a6226a6fa32`)
- go@go1.26.10/src/internal/syscall/unix/siginfo_linux_other.go (source SHA-256: `e2aba05f88bd7ca0e0fd1bf26fa5d53617398744dae3562bd4cee9937fd20486`)
- go@go1.26.10/src/internal/trace/tracev2/events.go (source SHA-256: `1f67012ee9a2b913027a173cf6ca34423cc33d9403536a7252cc9457d7894657`)
- go@go1.26.10/src/internal/trace/tracev2/spec.go (source SHA-256: `cf99c4a8b18901f580b3e221fde6c912519a7fd84025a7d3f42991c42c9dfb01`)
- go@go1.26.10/src/io/fs/format.go (source SHA-256: `8f0f84ae375932ef447ff0c3249548f7e0a3a7937bb0052341a521b78bd7d15c`)
- go@go1.26.10/src/io/fs/readlink.go (source SHA-256: `ffa574355de9871c697990db771088d92ad3f0ded7c6b9b701f9a6a60b81232b`)
- go@go1.26.10/src/iter/iter.go (source SHA-256: `ae73fb7933353c001ce3f85095ce248848a724e24490444918b467f725082019`)
- go@go1.26.10/src/log/internal/internal.go (source SHA-256: `8b55b6f684fde5bec66a1717b2475de06a050dca76aa0bdaf914a781e526eea0`)
- go@go1.26.10/src/log/slog/internal/ignorepc.go (source SHA-256: `ebc0bd29942e1313a93a7b52eb1578d11ae56d1d5e4a21a87513a5bd2643a5f7`)
- go@go1.26.10/src/math/rand/v2/chacha8.go (source SHA-256: `5ef2db8fe067f08142bd5867b0f7a27e5b2fd0f20178971d766d1fe4324460c6`)
- go@go1.26.10/src/math/rand/v2/pcg.go (source SHA-256: `6ebd291420193e3d127c3c5bc46e7b3ee99eb5db6c3729e3ac9786c912a0ffa5`)
- go@go1.26.10/src/mime/multipart/readmimeheader.go (source SHA-256: `395541411dd83b97939818f37053b634e9e37cf7b670b8ec2a368375922eeeb6`)
- go@go1.26.10/src/net/http/mapping.go (source SHA-256: `97be152228bf77d856b8ccf6f79cb41db42e244f63c2749905cd761ff1f943c8`)
- go@go1.26.10/src/net/http/pattern.go (source SHA-256: `e98488770e5e47b589add6245126faaec5e57bb0d2426ddeca774c889cf9f8ad`)
- go@go1.26.10/src/net/http/routing_index.go (source SHA-256: `05337a7df8a5c0e09004118326cfcaa89aa3ebf09f3aa8fb69727d8814dc3ef9`)
- go@go1.26.10/src/net/http/routing_tree.go (source SHA-256: `555d7ae16166612d9281b3f6f741c8e0a98f4e3a9968ddaea0b309182096879a`)
- go@go1.26.10/src/net/mptcpsock_linux.go (source SHA-256: `772be461c9f43ca5c647573a6b5ce595dc640206009ac41298c3e1106d52e55b`)
- go@go1.26.10/src/net/mptcpsock_stub.go (source SHA-256: `14927298389b19e73fb27cde42d10d9ecd56631e98bc5eb094c8417e29f7b982`)
- go@go1.26.10/src/net/netcgo_off.go (source SHA-256: `55966058ba76531fdc7f7e9216785b978387d7be093fc9c9176dc43494936eeb`)
- go@go1.26.10/src/net/netgo_off.go (source SHA-256: `2168ca7864525e92f205def9168355951f193aa2bdb415e6ad43c84d91f93c7a`)
- go@go1.26.10/src/net/rlimit_unix.go (source SHA-256: `bc6641769bf1262cc5b591785021807916a569bb9f0ece4562831e9a56206e21`)
- go@go1.26.10/src/net/tcpsock_unix.go (source SHA-256: `4761d925d681945c647092b78c0afc5ceb6a24e1c8c7ad2f44e631f2b9eaca2e`)
- go@go1.26.10/src/net/tcpsock_windows.go (source SHA-256: `91054c727154f7e7b3ad3dd44a19ebf0e8846b11df7b9bb79b30b6b887fb3762`)
- go@go1.26.10/src/os/file_open_unix.go (source SHA-256: `d9849e16c1714dd89b7787e2b5724d91e8053559f7bd6337ef22979a0782c238`)
- go@go1.26.10/src/os/pidfd_linux.go (source SHA-256: `f2a7c0e31bda943d4e2151a61f6a9eb2084815ed5c887cc8aa4fd7506805de33`)
- go@go1.26.10/src/os/pidfd_other.go (source SHA-256: `fad483ca74092ba62abc53c236e2b19f553e6dd3bbe5be11dd50c52d3f1fe751`)
- go@go1.26.10/src/runtime/coro.go (source SHA-256: `e510468afaaf06a14b515f6a51e2191fe1871c0655cbeef2c12268496a627c4a`)
- go@go1.26.10/src/runtime/fds_nonunix.go (source SHA-256: `f662698bfb09b33d344c46db013d1cf1b5ebef9a48e4b54a6a2735c199a3c6ee`)
- go@go1.26.10/src/runtime/fds_unix.go (source SHA-256: `82e2d0c6f73403b973c29c24e0e8aea8ad848f1c20c22cc1172a08ced04481de`)
- go@go1.26.10/src/runtime/minmax.go (source SHA-256: `3740bf87388f301a2a05c6d31e8d7160c9cf54dd970dcb365d4f8816e92ba084`)
- go@go1.26.10/src/runtime/os_unix.go (source SHA-256: `7a8d6fa105fa250c70b4cb9c90f46deef34d14449baaaafacb70c0f248067afe`)
- go@go1.26.10/src/runtime/pinner.go (source SHA-256: `6f638fb797edb1ed42601f11a086c8f6a34f095d7c85b9f53773e59528ff38bd`)
- go@go1.26.10/src/runtime/pprof/proto_darwin.go (source SHA-256: `956d80612ab4b40b27279ae6dcf719ece5e57d384af1aa7c08c252c8f2e3f8f9`)
- go@go1.26.10/src/runtime/pprof/vminfo_darwin.go (source SHA-256: `2883c97d83c2826a6e89eb42e8133505a03ae3ed744d26eba30476cdcc74046b`)
- go@go1.26.10/src/runtime/rand.go (source SHA-256: `f0c4935ece88ac3861ebd5060a3507831158bcf5ac47abae5860220c0da970c5`)
- go@go1.26.10/src/runtime/security_issetugid.go (source SHA-256: `ad7117b30a5a8870a64b938f2a7cec193617037af6c3cc414eb9b384cbf96950`)
- go@go1.26.10/src/runtime/security_linux.go (source SHA-256: `575dba66b2595e08b4cd7bb998a2f0d1040b3fb0c4410ce5d73b43f9e19dc2c1`)
- go@go1.26.10/src/runtime/security_nonunix.go (source SHA-256: `5026ff12daf3f61f487e47d7adf980d305ede84f3dbf66913b4603307cee64fb`)
- go@go1.26.10/src/runtime/security_unix.go (source SHA-256: `d570c7ee6b36cca9a0c270d641b2af4bd2b4e0f64e49dee7831019f04fc841d3`)
- go@go1.26.10/src/runtime/symtabinl.go (source SHA-256: `64bbd9ed4f94111659b89fc55c0764fb1f9a4202453ebb8646f7c6aef98c4edf`)
- go@go1.26.10/src/runtime/tagptr.go (source SHA-256: `4b46ea3f7fc1ea740de027c57ff42273b8b57ff76c692041ba43b8f78ee6caf2`)
- go@go1.26.10/src/runtime/test_amd64.go (source SHA-256: `432c8025a7a10b45e23f8e825244ca027e8dfc7cc123f1c8b58ea3a630923ead`)
- go@go1.26.10/src/runtime/test_stubs.go (source SHA-256: `be98082bf3f62a9ae7262ba89161c12b8a29328a83aa4976cd91c0d8e7187df9`)
- go@go1.26.10/src/runtime/trace.go (source SHA-256: `c4a410adcf39a2a822b1b268ad93df9b1515928161bf65a458e1bbbdd127bc51`)
- go@go1.26.10/src/runtime/tracebuf.go (source SHA-256: `ccd95094ef66367da61d82a82ffcb7b51ce537748666550811e650fc01f7c742`)
- go@go1.26.10/src/runtime/tracecpu.go (source SHA-256: `449de707b2c323a11f0aba7b9e3669ce18f0635247be664b01f3bf365156d670`)
- go@go1.26.10/src/runtime/traceevent.go (source SHA-256: `f8e91b49a351e77b59d28d5667a30a2df8e62c9270cf8eb47d0e983da17a4775`)
- go@go1.26.10/src/runtime/tracemap.go (source SHA-256: `41f2e3d3f6a34700495214e15e12f33da2500c61706161c593f39915e126118f`)
- go@go1.26.10/src/runtime/traceregion.go (source SHA-256: `2f6ae4e4c9668173d502ea57fa79b0d1ad5da8b154bc700bfbebcbfba8e17bf2`)
- go@go1.26.10/src/runtime/traceruntime.go (source SHA-256: `9cc5ffc359ac31558fdc2cf721ce70df38f9cfe0c28aef9755f71250a56ae0cb`)
- go@go1.26.10/src/runtime/tracestack.go (source SHA-256: `ce2e80a74992de3a2de6c23e77cc58009aff5cdf646b4e46b9cf1829e1fef63f`)
- go@go1.26.10/src/runtime/tracestatus.go (source SHA-256: `4db02b2b49c5bc682913d8265b0c4f4ce2913d3b2997c6722b1a9f1345140843`)
- go@go1.26.10/src/runtime/tracestring.go (source SHA-256: `fae042be62f6014265970e0977b1bd185575ca09f08241057deae3eb5c57caf5`)
- go@go1.26.10/src/runtime/tracetime.go (source SHA-256: `cfbb13de9ff44421db13f43ab22780f9da3a878982ea5671834b8ff0ee4c196f`)
- go@go1.26.10/src/runtime/tracetype.go (source SHA-256: `54a26cb246bb29e52b6fc267df3c0125ce1788816d543318ed1333708199313f`)
- go@go1.26.10/src/slices/sort.go (source SHA-256: `e6cd273b637e7dd2f547596c22e6236b92ae9b376f31ceb7b3c3269dd60a38df`)
- go@go1.26.10/src/sync/atomic/doc_64.go (source SHA-256: `0434f660a650bc6a1eaf9b8c48205ed0ce8c18b53001e949dde7662d04bde9ea`)
- go@go1.26.10/src/syscall/wtf8_windows.go (source SHA-256: `df35fc932a3bf98375eb780c08091f6a0b4a927ba2918a977495fe0b15c77f99`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/endian_little.go (source SHA-256: `c6bc70c372d9e1fe86fcf295f406b17bf04bf8d1af25c2456f58520cdaef3be9`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/runtime_auxv.go (source SHA-256: `d898ace395866bed261d403c8cd0ea6eab6d6d77f52042204957e389c38938cf`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/runtime_auxv_go121.go (source SHA-256: `6eee9d1a593dce53d22545dc5d5f6ed9127a43e6568ef5f5521920946510d445`)

````text
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/abi/stub.s (source SHA-256: `c24d5e3798ef3ae5d80f2b6049d30b2d5a7f5d6f531b093cfb1b6bfa8fbefc01`)
- go@go1.26.10/src/internal/chacha8rand/chacha8_amd64.s (source SHA-256: `ff54f61cfe85d12f67df5dd374635784a9899458242cfd867be325c9a0bdcbba`)
- go@go1.26.10/src/internal/chacha8rand/chacha8_arm64.s (source SHA-256: `21fbfcbae881a8117d968f7a5b3e52fb937570f1ebb3138ea8b0a6fbcbb7fcd5`)
- go@go1.26.10/src/runtime/cgo/abi_ppc64x.h (source SHA-256: `c805863eed3ec9d923afa76fb1948b6a3206cc55169ac9286537149671122d29`)
- go@go1.26.10/src/runtime/cgo/gcc_stack_darwin.c (source SHA-256: `406a0563a17fe501c6bfd5e41841c2d322caa4ab871a5bfe28d61bb36f678ed5`)
- go@go1.26.10/src/runtime/cgo/gcc_stack_unix.c (source SHA-256: `342b71a48decef1f94306fac076375cd0e0290cd34d134ddea09effe22a93877`)
- go@go1.26.10/src/runtime/cgo/gcc_stack_windows.c (source SHA-256: `9d801831c3beee705cb63081563adfe078f16e93c80c3cf53919a44f4d8e7d06`)
- go@go1.26.10/src/runtime/test_amd64.s (source SHA-256: `58ce05d3481714aec7b75abcf188137fb580c5ab671d62eedf4c9667b7d4ad25`)

````text
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/bytes/iter.go (source SHA-256: `628f8dc51af2c70325050c0125340ec01c6a7da84c53f2f5ae3ec040c7e5fc6f`)
- go@go1.26.10/src/crypto/cipher/gcm.go (source SHA-256: `8931e5cb0850b287f8910dd848019b7a6cc417f7be8556554f088afe29085c6f`)
- go@go1.26.10/src/crypto/fips140/fips140.go (source SHA-256: `8a9e62523e766e1717af8e9b1fbf6253ae81fa7ca0cad6e0f3beba131f55557b`)
- go@go1.26.10/src/crypto/hkdf/hkdf.go (source SHA-256: `7d9bc6828e075daecfaaabe5015e910fadc0a4264064b8b54962d07df778cb9d`)
- go@go1.26.10/src/crypto/hpke/hpke.go (source SHA-256: `48af513f3aca41e79583c63fd6d913213bbffb602a2060576bcdb1462e71c474`)
- go@go1.26.10/src/crypto/internal/fips140/aes/cast.go (source SHA-256: `7ac51de2b860827f13d33a3085449691dd23c31e620fd8a0a193ca8b20f00c8b`)
- go@go1.26.10/src/crypto/internal/fips140/aes/cbc.go (source SHA-256: `2539041bc7e23663d4f6953f40ce9751130db2d28d6f8548beca56636f35b5d5`)
- go@go1.26.10/src/crypto/internal/fips140/aes/cbc_noasm.go (source SHA-256: `d4bdaa6cd785d444157fd552748f5f5e53e4d13e456b0e7a29537781b96c1eec`)
- go@go1.26.10/src/crypto/internal/fips140/aes/ctr_asm.go (source SHA-256: `6e24fcf4398969616396fa9aeb4f5dbb6526fcc878f59516180f3f82ce4fec26`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/cast.go (source SHA-256: `78c62c041c1ce670b4e7d94ad47d2e0815f3270d563a4ee797f19f49b9e1d143`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/cmac.go (source SHA-256: `e0b6f3315d99d67b39b2f8e8b9f5c1040094503c2410c03f3cf5a43dba616990`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/ctrkdf.go (source SHA-256: `1e4258249d95d6e3d13f9df6549382b72b603565646f787a4c59c95efa4fa6f9`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/gcm_generic.go (source SHA-256: `def1fd72eb31bd265956d1c1cc7212a0debdb0349f6edef349a5901ab3576b00`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/gcm_nonces.go (source SHA-256: `7c31156b2fb9ab3ca6478332b52f68cccc086d361ab430a56eb127b48eb0ac6d`)
- go@go1.26.10/src/crypto/internal/fips140/aes/gcm/ghash.go (source SHA-256: `905a27c113c3837adbeb21f87e8acbe12abf2bfb5c018d0042930f691fdcba08`)
- go@go1.26.10/src/crypto/internal/fips140/cast.go (source SHA-256: `ed3b9fddc7581352d4fd67dae28569a33f965778330c9def418df7efc7f11e55`)
- go@go1.26.10/src/crypto/internal/fips140/check/check.go (source SHA-256: `31baa283c932141c381a22bbee2f589468926bd706fb04e59bc5254722390a0d`)
- go@go1.26.10/src/crypto/internal/fips140/drbg/cast.go (source SHA-256: `ad39245b1fb02b68daec596472085fb5c3a2e28f567a528a71278bc5cb084f92`)
- go@go1.26.10/src/crypto/internal/fips140/drbg/ctrdrbg.go (source SHA-256: `b9c7a312c89a1d91490ccc871dce6a3971044a87119978dfd0fd737a71b46990`)
- go@go1.26.10/src/crypto/internal/fips140/drbg/rand.go (source SHA-256: `51f50927a7e45df77ddacc2fd7b5206f25ed52e0e9896eaa6fa62b677514025f`)
- go@go1.26.10/src/crypto/internal/fips140/ecdh/cast.go (source SHA-256: `fa667658cbdf01dfdda1a96ea013f3aac28608c69a60aab0e793bae4b405bcbf`)
- go@go1.26.10/src/crypto/internal/fips140/ecdh/ecdh.go (source SHA-256: `1453adce8f850a6c1518f3a3969125812a57121d21004f8c4938d445f0b73007`)
- go@go1.26.10/src/crypto/internal/fips140/ecdsa/cast.go (source SHA-256: `68a18dc9cf665e03cc27be266a742efee60a2c2ee8dce29f156a4e5bd35e9636`)
- go@go1.26.10/src/crypto/internal/fips140/ecdsa/ecdsa.go (source SHA-256: `a8cada6e2681439e415b81b698f94aa499d1ba7157fc4fb9e34881ba9d082ec3`)
- go@go1.26.10/src/crypto/internal/fips140/ecdsa/hmacdrbg.go (source SHA-256: `13a4346d6df99e4ee7aae701b3e9dd737d4b3ed6d646758358cd161015378946`)
- go@go1.26.10/src/crypto/internal/fips140/ed25519/cast.go (source SHA-256: `44215bca1253e5d23c0a5222cdecdc1b3d7c3099d882cc2bfad97e2f0830efd5`)
- go@go1.26.10/src/crypto/internal/fips140/ed25519/ed25519.go (source SHA-256: `a9f71d91f3fc0228524b3ea58208e89937f34cd68eb7ae2e7933664104399550`)
- go@go1.26.10/src/crypto/internal/fips140/fips140.go (source SHA-256: `3428ba28d5aa395cdd3a9b01897b38cc017392e02b4bdfacc13e3e38b466ad11`)
- go@go1.26.10/src/crypto/internal/fips140/hkdf/cast.go (source SHA-256: `6579c46eb9886a6f423bc3b55ff8564c2d46aeb1b5358da4a5d00dad140acef6`)
- go@go1.26.10/src/crypto/internal/fips140/hkdf/hkdf.go (source SHA-256: `557d7be9d86e8112188d975c9b8376b468a937ca882e73090cf5f34d04722900`)
- go@go1.26.10/src/crypto/internal/fips140/hmac/cast.go (source SHA-256: `7fc2170a9524b86a66f66ef4be5f77d58f5e5f237de4c703599096d1e23a9bb9`)
- go@go1.26.10/src/crypto/internal/fips140/indicator.go (source SHA-256: `ac100029d7762a3281f76275f88b284c670fa097b3e86ff1860a3cbf73764c6c`)
- go@go1.26.10/src/crypto/internal/fips140/mlkem/cast.go (source SHA-256: `2a7ab7c4267861f48ed89bcddf5fac23bb5fe1d00e58b07d0e57ee1c5bad833b`)
- go@go1.26.10/src/crypto/internal/fips140/mlkem/field.go (source SHA-256: `2ef607225e367008c7eb5e4567ca8026fd39321c609d2aa95ce46998042c94a1`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/fiat/cast.go (source SHA-256: `7c461184661cb2ea99dba382796f4abc88be6544c77ebb4ee6f07d7058421403`)
- go@go1.26.10/src/crypto/internal/fips140/nistec/p256_table.go (source SHA-256: `ba494ee6fc933a189b5dad8cff8d86b8b19786cfdb27929453a4c851cfdc802b`)
- go@go1.26.10/src/crypto/internal/fips140/notasan.go (source SHA-256: `5d27c44171825d9c4597da23565e8e02a4eb5921609c010e5b24faf895dbab06`)
- go@go1.26.10/src/crypto/internal/fips140/notboring.go (source SHA-256: `36e73c6469f6f2f6aa5c9a5d6a789330cd53fd91365d011b305c1753da449588`)
- go@go1.26.10/src/crypto/internal/fips140/rsa/cast.go (source SHA-256: `a67fc5342bc97a553a442b57aac38b254e865b02cc32c411b5a530ec3e94ace7`)
- go@go1.26.10/src/crypto/internal/fips140/rsa/keygen.go (source SHA-256: `06aaf0880da2043862892261dd9773f1469537838f4e1b6b36af28dc71e7cf6d`)
- go@go1.26.10/src/crypto/internal/fips140/rsa/rsa.go (source SHA-256: `a3b73965247c692a6df2d5ce1d2fbf39b99121cc639d89a880f7fa306c7c2e95`)
- go@go1.26.10/src/crypto/internal/fips140/sha256/cast.go (source SHA-256: `69d07cca38c33f57bd7f6baa36a4d0c69b21aa4d6a1c1ed9d32dedd196b97b2f`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/cast.go (source SHA-256: `4fd14d4f9c5265e93ca01389533ff21029a7fa9f6cd4d354db099afa1c591c48`)
- go@go1.26.10/src/crypto/internal/fips140/sha512/cast.go (source SHA-256: `636a680a20fab95604ac603e337230c7f5ebc06563dda40aa2fdc13d8235c55f`)
- go@go1.26.10/src/crypto/internal/fips140/tls12/cast.go (source SHA-256: `55f9a78268d7dc14e9fc19a9606ccd7c52ab6a5baa2cf79f0051ccce1db24531`)
- go@go1.26.10/src/crypto/internal/fips140/tls12/tls12.go (source SHA-256: `5fa5cac46973f0cd34652ab1535213adc1c8e7cb51041e52a0926eec96d601ad`)
- go@go1.26.10/src/crypto/internal/fips140/tls13/cast.go (source SHA-256: `162c9914d09f41d0f5bf4ee99c1bd5ec89ccb4c0bec1ecb3c6ffa90109ca302f`)
- go@go1.26.10/src/crypto/internal/fips140/tls13/tls13.go (source SHA-256: `c04b704715d4fab230ed7257eb46f4414c9eb30d99944947b1fe68fa6dded433`)
- go@go1.26.10/src/crypto/internal/fips140deps/byteorder/byteorder.go (source SHA-256: `8ac6a9915fdb1ad6d05e8e40e0dd694102bcc529a7e79acb8c7859a1c70caf38`)
- go@go1.26.10/src/crypto/internal/fips140deps/cpu/cpu.go (source SHA-256: `0a17476e7b34fba9e513e3bc9a34f300185ff39f44c45e4b309031d4e2d1d453`)
- go@go1.26.10/src/crypto/internal/fips140deps/godebug/godebug.go (source SHA-256: `dfbc907e235aea8da676e585cc409d455866b2b70161668a6817ccf296672d46`)
- go@go1.26.10/src/crypto/internal/fips140hash/hash.go (source SHA-256: `860c04522dee5d560659c7f54a3d96b3aebd0077728f89c88a97c882cfe33ab3`)
- go@go1.26.10/src/crypto/internal/fips140only/fips140only.go (source SHA-256: `9857dbf233414f443af96a6421946c69192360f0e41fe9e13f5580bb9ab580f6`)
- go@go1.26.10/src/crypto/internal/impl/impl.go (source SHA-256: `0ee5bea6ebb52253d11efba0741b4bcd43c2172ddc328fd17f63d0e2da999567`)
- go@go1.26.10/src/crypto/internal/sysrand/rand_arc4random.go (source SHA-256: `62545bde8c7adc3c0c415ae8b06d0bfbdf8a9ebb3866a64498b9f67e4ea2257e`)
- go@go1.26.10/src/crypto/rand/text.go (source SHA-256: `c2c7d3bc9f093676f362b9e22772b6ab2d48acd6c06ca85fea0fef4741c17ce4`)
- go@go1.26.10/src/crypto/rsa/fips.go (source SHA-256: `380445d315516f986d5bc28e100b3d227acf121b8d70e04dcc28918c782d981d`)
- go@go1.26.10/src/crypto/sha3/sha3.go (source SHA-256: `b0ab1d099bc7779180dec6ef1c18ecb8dde324c81667d260dfa010e89db72638`)
- go@go1.26.10/src/crypto/subtle/dit.go (source SHA-256: `6049fc374092a432b4af6693c898905723d035c691e6a9dbd752f91557ff5aa6`)
- go@go1.26.10/src/crypto/tls/defaults.go (source SHA-256: `937206b6ad70e36e9ddc8aeb036091e948f49bd07e8d7280580627a95c6edb39`)
- go@go1.26.10/src/crypto/tls/ech.go (source SHA-256: `8f6c8035a9928fd1d7b25be1490697618f493e9bbb0fb88a990fd6e5e0a66d7a`)
- go@go1.26.10/src/crypto/tls/internal/fips140tls/fipstls.go (source SHA-256: `19ab648f8b24078cc220612bb02adb8617b0a2ea327636707ca71343b2cc423e`)
- go@go1.26.10/src/internal/abi/escape.go (source SHA-256: `494d85e5231a20cd88b56fbcf02028e927a56e62c8d4f7212c08156bbed58457`)
- go@go1.26.10/src/internal/abi/rangefuncconsts.go (source SHA-256: `b0350d4021e6e279675e082416cc9280c296cfff6a123863ed7f2d693c82a14f`)
- go@go1.26.10/src/internal/abi/runtime.go (source SHA-256: `4e5e32bebcb05a34756123a330e4c4a79c0e4eb2fba3e81253fd3e616aca50ec`)
- go@go1.26.10/src/internal/asan/doc.go (source SHA-256: `dfd59da21a8923944fe984eb6d64ae87ccc9eb3b9bc895240c20ec7db6115e81`)
- go@go1.26.10/src/internal/asan/noasan.go (source SHA-256: `7de9002b9b667f10f43cab22616ad65348a8e35d0a1c29ef871f817c27211143`)
- go@go1.26.10/src/internal/byteorder/byteorder.go (source SHA-256: `c31dade373ba9a6970fdec0943cdf95b2a8ae1f0e057fdd581e50165aca26fbb`)
- go@go1.26.10/src/internal/filepathlite/path.go (source SHA-256: `dd1a0138328d6353d089d6ba8e04905d9e951bedfef1b85913c64de2ca2fd223`)
- go@go1.26.10/src/internal/msan/doc.go (source SHA-256: `eed9a5e6933f2adc99ebc4486a0e22bc044dc4ca71c4f3df3d2cba3131f70fad`)
- go@go1.26.10/src/internal/msan/nomsan.go (source SHA-256: `a93262f02003de244b1490022bd6e8765c0efad5a25372450bf87f9fbedb8470`)
- go@go1.26.10/src/internal/poll/copy_file_range_unix.go (source SHA-256: `ba0ca485fc1c951a4011dfb7a760c6fdef8d9911b5cac0696079381f62d37edd`)
- go@go1.26.10/src/internal/poll/sendfile.go (source SHA-256: `86cbbbbe929c151795ac01e6b939e72838f072fefaa1959887badc5d531bf955`)
- go@go1.26.10/src/internal/poll/sendfile_unix.go (source SHA-256: `912a2ad7eb419472c6348689f259d6d81e124cb82593789e03741d078ff66be5`)
- go@go1.26.10/src/internal/profilerecord/profilerecord.go (source SHA-256: `3d866406f890f806fbd2a3473fa1557067a5963d3e072910e1be7d48cb69b6ac`)
- go@go1.26.10/src/internal/runtime/exithook/hooks.go (source SHA-256: `1ee8780b40bd8512ca2bbe16a50b9a4f0ab94e4362bd4c2017e0ac3446880243`)
- go@go1.26.10/src/internal/runtime/maps/group.go (source SHA-256: `cd99fe64dc1c3aa3a80516a070e6826afc424ddb11e883cd9dd761a2d5570aef`)
- go@go1.26.10/src/internal/runtime/maps/map.go (source SHA-256: `cec1f588f261caf8fe4a9b310fae11abcc162d5c975be21d243f9ce727151432`)
- go@go1.26.10/src/internal/runtime/maps/runtime.go (source SHA-256: `51e7dc517a8f14e85c1da24ffd72fda950280fa39e0f53cfafded0dc447cce86`)
- go@go1.26.10/src/internal/runtime/maps/runtime_fast32.go (source SHA-256: `99cc80f2d14739205e12b8096b02a0445f4641b3caed6e697d09d41bd9d5af08`)
- go@go1.26.10/src/internal/runtime/maps/runtime_fast64.go (source SHA-256: `29549c9d433050cfba341e169137f5f0b5748c5c43936f7dd44f38609117b52b`)
- go@go1.26.10/src/internal/runtime/maps/runtime_faststr.go (source SHA-256: `9a4abcef59b9a68f01a20da7d68b783c0b3c276a0b876810e45cdc5b4965e1ab`)
- go@go1.26.10/src/internal/runtime/maps/table.go (source SHA-256: `d15567aa1735a8221599af3b11f857292a7b270daaa2820ff4e8053d4f172112`)
- go@go1.26.10/src/internal/runtime/maps/table_debug.go (source SHA-256: `6e693977c479085d77611f3c3b482e5fb27146c8fa1dbbb299482ad2eb11430d`)
- go@go1.26.10/src/internal/runtime/sys/dit_arm64.go (source SHA-256: `ae8421bf5674c0bf2dd028da42c93212032350626d5074d75760fbac883cd177`)
- go@go1.26.10/src/internal/runtime/sys/no_dit.go (source SHA-256: `144597c72d957b854ca514bd0f1e47893f5fb76c62fabc2e95162c32067e90d4`)
- go@go1.26.10/src/internal/runtime/syscall/linux/defs_linux.go (source SHA-256: `29b7f85c166cdb36e21b60ba259ab76cc560bb67dfaf452656393d90190a1d6e`)
- go@go1.26.10/src/internal/stringslite/strings.go (source SHA-256: `764f3227e30300d0227b8c971aa718032232aa16370535476dd8e7e830fa4985`)
- go@go1.26.10/src/internal/sync/hashtriemap.go (source SHA-256: `ac36cd5db5c5c659b17255056d88b25c132ed1550bc35c1437be70e0e7d537cc`)
- go@go1.26.10/src/internal/sync/mutex.go (source SHA-256: `0535ca03e44cb62cf3e4a511b989ffd0f08ec899d864e2a817b77c4025bec48c`)
- go@go1.26.10/src/internal/sync/runtime.go (source SHA-256: `e3445cf3a17f6b322cb3136e30a28eab7a70c0f37fac44ba5dc7912a19930a17`)
- go@go1.26.10/src/internal/synctest/synctest.go (source SHA-256: `3291fd3debb2430be3017c9a3cd783ad4d1a2791d8b4133b473a11a29e250449`)
- go@go1.26.10/src/internal/syscall/unix/arc4random_darwin.go (source SHA-256: `e25578e3c171f0712232f6ce4dcc98f62c7121a4aa3b90a2046832ec2ab19eb7`)
- go@go1.26.10/src/internal/syscall/unix/at_darwin.go (source SHA-256: `eaea1b0e18d58dc7f18adff2b1915db78bb0eaa7ae486f5424f0057355ec7e84`)
- go@go1.26.10/src/internal/syscall/unix/eaccess.go (source SHA-256: `f81216e0796d54c28e4e33328c0c5b75fe98b7dab172ac30d35be304412c9b97`)
- go@go1.26.10/src/internal/syscall/unix/faccessat_darwin.go (source SHA-256: `3765abe46d2ecac27a64d04578af9cb6f8c3a22c6d5a3935526253b825c155fa`)
- go@go1.26.10/src/internal/syscall/unix/faccessat_syscall.go (source SHA-256: `60d3fb3116514098e471e948e845b77582d2e777859e3dde03cf73d65b99fe69`)
- go@go1.26.10/src/internal/syscall/unix/nofollow_posix.go (source SHA-256: `5c5562ab0d8f94c8108172ea142ce5d0e20a07e471e0c7bd9d9bede2c6c93689`)
- go@go1.26.10/src/internal/syscall/unix/syscall.go (source SHA-256: `65de2b30e89644981e8dc159d226ebe585a66112438e382e6cd2651aa6f68451`)
- go@go1.26.10/src/internal/syscall/unix/tcsetpgrp_bsd.go (source SHA-256: `594ae40596c382d4593abec994d60c058dcd2ba6f897487a4f4a1268421271b0`)
- go@go1.26.10/src/internal/syscall/unix/tcsetpgrp_linux.go (source SHA-256: `a154792f653632733e8a2c81fd2620b0fa08b0e977774af6c64beaefa970855c`)
- go@go1.26.10/src/internal/syscall/unix/waitid_linux.go (source SHA-256: `6e59316cada5dc52c568b2bc62e7f2220a72114b3438beb11dd6d8f7abdaf68c`)
- go@go1.26.10/src/internal/syscall/windows/at_windows.go (source SHA-256: `397cfa7c20513134897c3dd6c265d2f88e2daa7577e985d23245d06cf6715f75`)
- go@go1.26.10/src/internal/syscall/windows/string_windows.go (source SHA-256: `946542ece8b41098ee14238b73402975fb8271878f1a5694527e95636d27dbaf`)
- go@go1.26.10/src/internal/syscall/windows/types_windows.go (source SHA-256: `07b223c29af0acb0b8251415e362635168c657ba6caa8dd3544e97c5fd489a80`)
- go@go1.26.10/src/internal/syscall/windows/version_windows.go (source SHA-256: `fde689736af5573dc20621ba662ebd2fcbd22780b75bbfb081ec6a169fddadd4`)
- go@go1.26.10/src/maps/iter.go (source SHA-256: `a803b5b1b09aa2dcedf90db32270f9cec310bc7cf8f141f87c2e6567db0c4208`)
- go@go1.26.10/src/os/eloop_other.go (source SHA-256: `c39ebc651b2284508b70c299c806e9becca82515b2ee4fe38c918bec7a046ba8`)
- go@go1.26.10/src/os/exec_linux.go (source SHA-256: `b80854e3f1f6ce249d9d04d9198d4081f8cfb67845e49a57565fa1118e279dcc`)
- go@go1.26.10/src/os/exec_nohandle.go (source SHA-256: `4fc2b84936dc0dcbc1ce33db663803947d2f76238c295fbeccd92219588c0a15`)
- go@go1.26.10/src/os/root.go (source SHA-256: `56355baf52dd718903a59c94d54eee5e2765d50383b471d08819accd1ec02acd`)
- go@go1.26.10/src/os/root_nonwindows.go (source SHA-256: `49152bef7df7c862984d95000be420f4f73afc56ef1e5abb7caeb6e8b5e1f5de`)
- go@go1.26.10/src/os/root_openat.go (source SHA-256: `ed015476176e245d77cd330242eb19fb4abcb9b81b125eece3b43b3a290f35f2`)
- go@go1.26.10/src/os/root_unix.go (source SHA-256: `a4ab2861f9f251650b70d4a27a9a1fb3ac4a10676369862f2ffac3ca3ecdbb0a`)
- go@go1.26.10/src/os/root_windows.go (source SHA-256: `2c8ddba33817157b0dc9c5b6e2f5ab7aa0285b5722d2d63b8ac5727f253aad9b`)
- go@go1.26.10/src/os/zero_copy_posix.go (source SHA-256: `3e2916cf88d98c98a4b4810754fbced3897ba9f6d58179edc88f000a10987ab6`)
- go@go1.26.10/src/path/filepath/path_windows.go (source SHA-256: `96c1af1caace59a146241f5a48a700a142f2e5520f932fb2fc6376ce5a94ed81`)
- go@go1.26.10/src/reflect/badlinkname.go (source SHA-256: `51e6f422e442cd7c091d5c67e47b505ba7c59b0cd26017246d084221decad82c`)
- go@go1.26.10/src/reflect/iter.go (source SHA-256: `78128c2b362102626bcdd579a553543f3e5a7be6ab6a02d04e0662a3fcbd71f0`)
- go@go1.26.10/src/reflect/map.go (source SHA-256: `48f44acf853d0b06fe3ae5af4955934d736fda4926026c68b8faa073c7d2b8d4`)
- go@go1.26.10/src/runtime/badlinkname.go (source SHA-256: `172c02155463385704c6f873e8d05e4bb3a04054cca508b158c203b54ee205ea`)
- go@go1.26.10/src/runtime/badlinkname_linux.go (source SHA-256: `3defa450ca9d61ab27d61d28308c36e35a554f759f4e3fd444632aba37b8c0ae`)
- go@go1.26.10/src/runtime/linkname.go (source SHA-256: `4165304e269d11107daef57cb81cb1e7fdb2e5a061d882330d42cb5a9c84138f`)
- go@go1.26.10/src/runtime/linkname_unix.go (source SHA-256: `e77e8fc214ad311a5a74b48585d60ab46a91514049e6c66697bc1e873ced9922`)
- go@go1.26.10/src/runtime/lock_spinbit.go (source SHA-256: `fd134af1468b712b419c6be01b37d7c5b065cca043ca1d30350c73ce3e6c1821`)
- go@go1.26.10/src/runtime/mcleanup.go (source SHA-256: `bda3732f12dc16750994d6ce6b10b6844beb2ea1aa854c3407be757097f04207`)
- go@go1.26.10/src/runtime/mem_nonsbrk.go (source SHA-256: `975fb2b84c6f9cba872fd6e3f77726ea297959a7373d75794a58194860ac07e7`)
- go@go1.26.10/src/runtime/netpoll_kqueue_event.go (source SHA-256: `63d8d7d0642dcfaaa0e4b9531e48a1f3a8393009946f2152654e96b5f11bed71`)
- go@go1.26.10/src/runtime/note_other.go (source SHA-256: `8f910c2a6c624eba7f47275cb2917de934ead4300b5de64bea5637b5de1ee7a2`)
- go@go1.26.10/src/runtime/secret.go (source SHA-256: `a9f057fc28b7bc748f689b54246b97eabaaec395cba732da721bc103be848c42`)
- go@go1.26.10/src/runtime/secret_asm.go (source SHA-256: `a946d8499d78b76ae8ebe3e7082c511602518501bddf924449ab8889c8c8276e`)
- go@go1.26.10/src/runtime/stubs_nonwasm.go (source SHA-256: `b5e6633be53365857e330e367fede63060e8ec23139850f2a643a63972873034`)
- go@go1.26.10/src/runtime/synctest.go (source SHA-256: `7c9fe55993868e9b06502906b38cf6d72b763b224d0f0ef605e863b4d39f4272`)
- go@go1.26.10/src/runtime/traceallocfree.go (source SHA-256: `096a01883b7179bc9789b2f92bcca9b6314641c2406cc42e3333852adb32ef49`)
- go@go1.26.10/src/runtime/vgetrandom_linux.go (source SHA-256: `fda3cf62391408d701c497a0dbbea92ae7ebddc99347a9efa9e0915afce90fd2`)
- go@go1.26.10/src/runtime/vgetrandom_unsupported.go (source SHA-256: `3c6f633f35724bf61203fe3abb4047a0c0544fb1dd86f7f698641f9b57cf8f55`)
- go@go1.26.10/src/slices/iter.go (source SHA-256: `986cce7d58c65e6f87987544908ec138ae664c4bf5c8868d2177a425f6fa2bbc`)
- go@go1.26.10/src/strings/iter.go (source SHA-256: `067603fe06f2cec8ec607f9c9d469493901820f76af7f9e3e51136824fe8b04a`)
- go@go1.26.10/src/structs/doc.go (source SHA-256: `f9586c696223039cacc8c12753088e5c95188252abda95015566246e1138f81e`)
- go@go1.26.10/src/structs/hostlayout.go (source SHA-256: `322d314ded4ace4a48b0fcc7f3e1979bde678a69195f835d111e9ba38a5cb0e7`)
- go@go1.26.10/src/sync/map.go (source SHA-256: `ea2b7b418426ba83ed9b77c5926dcb28a6a94dc84b26ea7b581c782a36a646b4`)
- go@go1.26.10/src/syscall/badlinkname_unix.go (source SHA-256: `2d5d9bc9885f05f58ce7f7334232ce15668931aadb9eb6c49c71b8a43fa3f83e`)
- go@go1.26.10/src/syscall/linkname_bsd.go (source SHA-256: `5ea0e17dd9063bbd644cb4dda35128f913480d1a978d84d03a7264d028619cde`)
- go@go1.26.10/src/syscall/linkname_darwin.go (source SHA-256: `08afe9bbb20b068de8f77b67a2fa1c565b75f6063f25bc05de31a56dbd0611bf`)
- go@go1.26.10/src/syscall/linkname_libc.go (source SHA-256: `37d821f5c215d7560d4acccb73ab712bc165c43109b66273206aadb95a1ba40a`)
- go@go1.26.10/src/syscall/linkname_unix.go (source SHA-256: `5bf0a04177de7a0537a305c39cc9e380f626a04c569b88ac3679d9acdcef2924`)
- go@go1.26.10/src/unique/clone.go (source SHA-256: `ca3a74b2f8a1fb1dde53b6a539d4ca0868acf668e4c75bec8759074b1af728e8`)
- go@go1.26.10/src/unique/doc.go (source SHA-256: `62e40bcdc338d51d33429f67f2393d3282a75285c0ae2ae6b380c8e888e05dca`)
- go@go1.26.10/src/unique/handle.go (source SHA-256: `08c7ca3c26f66f22128d0f061e47242fc3a96e6dd765ac73d1fb8800e022c859`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu_darwin_x86.go (source SHA-256: `4d90cd4b64eed4f9bbfe8e6c4b856f785659d6f22f56f8f6cf1b172c10e9c518`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/cpu_other_x86.go (source SHA-256: `6f447bdafdc6f75593b29977c4f43a39f743043d4f286c542a8b93fa7fad3bf9`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/syscall_darwin_x86_gc.go (source SHA-256: `66059b45534544bdcde771fa793eca91d04de079e21adc4b2732e0912a57900a`)
- go@go1.26.10/src/weak/doc.go (source SHA-256: `2446a5aca7a3859e805750097c52a99ec97f58c4dc117700d8cc835af8100352`)
- go@go1.26.10/src/weak/pointer.go (source SHA-256: `3c4efb4f12ad6a38adc33c63a8835f062051cb789b7973619baaa18451b65520`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/runtime/sys/dit_arm64.s (source SHA-256: `31bfeabfbbfd893210a062be1f722800e5e290a4d365d376bf6b4887be9002f5`)
- go@go1.26.10/src/runtime/ints.s (source SHA-256: `b7a34cb11c9d0fe5ff9ed575b74b2716f3e86fc44c4044589770baf93b52c867`)
- go@go1.26.10/src/runtime/secret_amd64.s (source SHA-256: `9b61f272835c56f61f26bd9b1e3f5f976c779dc5ada377cadb3742c68c34071e`)
- go@go1.26.10/src/runtime/secret_arm64.s (source SHA-256: `12938eddaf0ef632af8ae3e3d4cab316f271b5845793f2a172259268613f1a17`)
- go@go1.26.10/src/vendor/golang.org/x/sys/cpu/asm_darwin_x86_gc.s (source SHA-256: `21022bf7461861b3744775075f4a4ef960c036cc469914684f87a4da9fc8b3a8`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/fips140/enforcement.go (source SHA-256: `b8106111dbbf1a1b2eb3b7123d5d6affe61cc22b4f5ea4510fc187ac0e401531`)
- go@go1.26.10/src/crypto/hpke/aead.go (source SHA-256: `c21e5faf4ec6002345c68c4964b194d1df4189174b59a503a63509d48f671975`)
- go@go1.26.10/src/crypto/hpke/aead_fips140v1.26.go (source SHA-256: `1a3ae9c24fce2e10ddd51d274eb4d44708c40c78e3537de971ce48c804a8262a`)
- go@go1.26.10/src/crypto/hpke/kdf.go (source SHA-256: `54c5659c70a27d6553fae6233fa0db753a674b7e771dc62e9e714f4f0d99f940`)
- go@go1.26.10/src/crypto/hpke/kem.go (source SHA-256: `1f1414cf5951926fd7bb410ac7d0df84b49b042ee838c5678e06e745c0ea0e3b`)
- go@go1.26.10/src/crypto/hpke/pq.go (source SHA-256: `e19b503aad3cf964bc57547705f2638fbee25704bfe318a726f584776a4a5906`)
- go@go1.26.10/src/crypto/internal/constanttime/constant_time.go (source SHA-256: `ecc01311f2c0e9ffd370180a9d65fba21d043417e7408ed924327ebc46f8ef95`)
- go@go1.26.10/src/crypto/internal/entropy/v1.0.0/entropy.go (source SHA-256: `630e9ba82e5b1282c9976de6fedbd784338e5aea20fb4b0ffb7c36604b3495ff`)
- go@go1.26.10/src/crypto/internal/entropy/v1.0.0/sha384.go (source SHA-256: `8e5d6d543835c161ce50031c711fab4229c3822f29ba1967cf81da1bd73acf39`)
- go@go1.26.10/src/crypto/internal/fips140/notpurego.go (source SHA-256: `9daf0d4bbaca8654bfe6cc222987312f497e535d23b3c2df4b6ed7f79701948c`)
- go@go1.26.10/src/crypto/internal/fips140/sha3/sha3_arm64.go (source SHA-256: `0edfd39f3c3a23a0127bac6da3067047d9cc54e1074476288ab4147b024871c3`)
- go@go1.26.10/src/crypto/internal/fips140cache/cache.go (source SHA-256: `d2052e4809787725708ebcbf88d10747e1031e1876fb46189ea1c2c3fd29b519`)
- go@go1.26.10/src/crypto/internal/fips140deps/time/time.go (source SHA-256: `7b30905ad2c20194d4818ca2bfa75fc44a2239f5703a2e208ab97bfc955f413d`)
- go@go1.26.10/src/crypto/internal/fips140deps/time/time_windows.go (source SHA-256: `ae818fa65eb74cf835fa588cbe265146254125c17899dcf10cf630221ffd304e`)
- go@go1.26.10/src/crypto/internal/rand/rand.go (source SHA-256: `fff364e2098359a767f47041c9e2939ea76850424a36183adfe748c4ce5c6da3`)
- go@go1.26.10/src/crypto/internal/rand/rand_fips140v1.26.go (source SHA-256: `0111d4ffcc48806c868ae1f3b4449e42e6af62ea21111bff1b9c3246b8ecc5b9`)
- go@go1.26.10/src/crypto/tls/defaults_fips140.go (source SHA-256: `0c8bd50bb3e2d3f7d46d905dc98b1845840fdd41e4957ec6aed902b9c6a9004b`)
- go@go1.26.10/src/crypto/x509/constraints.go (source SHA-256: `7a028a170f8371be1eee3311978b7bd88eade5246a864d4537c3c7cd39bb51bc`)
- go@go1.26.10/src/go/token/tree.go (source SHA-256: `ef1f838eb2a9e3d83c319b6bf4a1e8b8c32fdf73aab42c315c75308ede05abc5`)
- go@go1.26.10/src/image/jpeg/dct.go (source SHA-256: `c16e9585a6fb9ba08a0aacfd36911bf394ec3abaaaa633a2589256101be16d30`)
- go@go1.26.10/src/internal/abi/bounds.go (source SHA-256: `ea44135b94ceb854b298437fd4e6766dc7ebfeee51f9d1c9dae9cb273cc17134`)
- go@go1.26.10/src/internal/cpu/cpu_x86_darwin.go (source SHA-256: `036050c7259aa9282697b31603c3e5c8eb1ce234a1bbf26f9a326635a9645040`)
- go@go1.26.10/src/internal/cpu/cpu_x86_other.go (source SHA-256: `f24c84fc53bb2b8b067fdb462d6ac56594c371163527fba998f45c0f65c0c2f9`)
- go@go1.26.10/src/internal/cpu/datacache_unsupported.go (source SHA-256: `1ee1b485113109cb851e7141c072ad31493414c963fd469f398fc0ea2b385a9c`)
- go@go1.26.10/src/internal/cpu/datacache_x86.go (source SHA-256: `bbd50f64070d9497f062de27710a17d1aead867054d7dd9143bf063cdcd4fbea`)
- go@go1.26.10/src/internal/runtime/atomic/xchg8.go (source SHA-256: `9ebd5173979bd1c231f8a1c807820a83c27064ce03c0b57d732611d5e4fb8c7e`)
- go@go1.26.10/src/internal/runtime/cgroup/cgroup.go (source SHA-256: `e928d9fa44e1d758cfa38a7f4f140bca8c3e2698dd6e1f79b2eeecf40c5c94d4`)
- go@go1.26.10/src/internal/runtime/cgroup/cgroup_linux.go (source SHA-256: `9f192b16192c9176b5c1bc100e786f5c726359d86782aa86a32221bebb07d031`)
- go@go1.26.10/src/internal/runtime/cgroup/line_reader.go (source SHA-256: `14ae28a739829f8e3139fd77152bd9c94a940a3c33102bda7eb43cdaa2406529`)
- go@go1.26.10/src/internal/runtime/cgroup/runtime.go (source SHA-256: `3bc58574f37162ffc7519b72bcfb81614a831547a2fc2d85cd3ea960763df597`)
- go@go1.26.10/src/internal/runtime/gc/malloc.go (source SHA-256: `889e4ee48345fd246b261a9352a0f8b9deb2c595ab59629c0ca70f7eaf06462c`)
- go@go1.26.10/src/internal/runtime/gc/scan.go (source SHA-256: `f3cb2baab1c1df77fedd8904d9def2f19c92d628111d9fb026b70d48690cc987`)
- go@go1.26.10/src/internal/runtime/gc/scan/expand_amd64.go (source SHA-256: `7e2a099fab42af62a4e7340425ffdb9cfb9b2a91b3886d29ee018e1ddddbbb84`)
- go@go1.26.10/src/internal/runtime/gc/scan/expand_reference.go (source SHA-256: `58e336c9e78a10fb5200451ca52a991b6c5663f224954af6b6ede0a6cec3efff`)
- go@go1.26.10/src/internal/runtime/gc/scan/filter.go (source SHA-256: `87b9607d74699084abe2df8be62f7eb97793132d6a1bb97cc99e851900372230`)
- go@go1.26.10/src/internal/runtime/gc/scan/filter_amd64.go (source SHA-256: `21e405b81d2200a049db1131aa6891d58cfd86aea41e3b3d2560dfa164e23e9c`)
- go@go1.26.10/src/internal/runtime/gc/scan/scan_amd64.go (source SHA-256: `274b655592508439929493836be2abbd6ae5011902582c53a7f31a3a11242017`)
- go@go1.26.10/src/internal/runtime/gc/scan/scan_generic.go (source SHA-256: `4a588096ec337b47439a252de9b7fadb65af7fa27d791d3f0b46c70ef77162b4`)
- go@go1.26.10/src/internal/runtime/gc/scan/scan_go.go (source SHA-256: `6cd00f9633cc841915f402fe0fbb96a377fb1860da9cd1da499f966d52bd3686`)
- go@go1.26.10/src/internal/runtime/gc/scan/scan_reference.go (source SHA-256: `41e1bc4eea7370d1d37989a0c87c441f56e27701beec979425bb5800c21bdcba`)
- go@go1.26.10/src/internal/runtime/pprof/label/labelset.go (source SHA-256: `248d89e463db994678c7fec529fc5e523dda82c579ba7305a51b056292e7df58`)
- go@go1.26.10/src/internal/runtime/syscall/windows/defs_windows.go (source SHA-256: `8a3471a8d0307af983275d6c1ce9b503cc22375affc95c7eef7e0c50ee9adde0`)
- go@go1.26.10/src/internal/runtime/syscall/windows/defs_windows_amd64.go (source SHA-256: `bc8e87e75eadda689c833f8353a4bdc3cad0310d173ca9617d53e13aa841d7b4`)
- go@go1.26.10/src/internal/runtime/syscall/windows/syscall_windows.go (source SHA-256: `8d51118789a846f0eed3e41b1747161f50fb98cdf468434eaced39dd176fcfbd`)
- go@go1.26.10/src/internal/strconv/deps.go (source SHA-256: `052459732bddd00d8584351fb6a906566fc664ea40cbb125a6084c266ec32238`)
- go@go1.26.10/src/internal/strconv/ftoadbox.go (source SHA-256: `131ede3bec4e39c5737c83fc074401bf22fa4e03b2038c972a039310f1f7dfb3`)
- go@go1.26.10/src/internal/strconv/ftoafixed.go (source SHA-256: `6b8a2f52d5e53206930ba481165268f58d08b5477a4810b91eca55345fecc093`)
- go@go1.26.10/src/internal/strconv/math.go (source SHA-256: `88db7fb281cdfab21f657eab8b37af131f23d822977c7e71bc753ffc5aec0145`)
- go@go1.26.10/src/internal/strconv/pow10tab.go (source SHA-256: `1c4356efc589a293b12b6af86f0fca7a88d139d01876ba98306cc1aae2e19ab4`)
- go@go1.26.10/src/internal/syscall/unix/kernel_version_ge.go (source SHA-256: `4732fa4207a03845d3f3ab21b93d951a04ea48fd10a77ee3dd41d0d55ff12a26`)
- go@go1.26.10/src/internal/syscall/unix/renameat_sysnum_linux.go (source SHA-256: `8c0c2f88d0c7e34b9caa1a21b771fa40d284b3844b02e63f194f26e2e550377d`)
- go@go1.26.10/src/internal/syscall/unix/utimes.go (source SHA-256: `e37b2cfa5ff403dff694dc69dfef3fff40bad5a6d99a311d1dc123c1ff469deb`)
- go@go1.26.10/src/internal/syscall/windows/nonblocking_windows.go (source SHA-256: `76016319aa6db9f7ca0ba9be706271dffe6dde24925196e97ae9f025ace18e81`)
- go@go1.26.10/src/internal/trace/tracev2/doc.go (source SHA-256: `021b396705a18015e433d99bf1883bc64ccfc838bb487b4133a423306cb4e075`)
- go@go1.26.10/src/log/slog/multi_handler.go (source SHA-256: `6e48fe54f7881d376d41d628f303adba6f60b6f2da826389db7523bc2ce8a414`)
- go@go1.26.10/src/net/file_posix.go (source SHA-256: `db3c6b88dad9f5ac472d2aebbac3e2ccfa18645718d99fde31a31967b8021bac`)
- go@go1.26.10/src/net/http/clientconn.go (source SHA-256: `afa8a5e6768d52c591901e4915de9669ff221f90013744e86819de4a5dc0d324`)
- go@go1.26.10/src/net/http/csrf.go (source SHA-256: `9376ac5c11e3205c51eb596691c18743269deb7e0bf5036023f874bb1bf4f3e6`)
- go@go1.26.10/src/net/sendfile_nonwindows.go (source SHA-256: `8e0502d6901abb50af70ff3b5e3e6b4bdebc809b883f1d2fc858aa118e1a3cf8`)
- go@go1.26.10/src/net/sendfile_windows.go (source SHA-256: `38c4d3cb46d29b187cbb97be77dd87202d09b5b7b37b665f20cfdca734ab9f3d`)
- go@go1.26.10/src/net/url/encoding_table.go (source SHA-256: `274e5ad203d6c38c41fe435eb9e5fc65ae1a38e98b44daa9b6355ba0d3fdac09`)
- go@go1.26.10/src/os/removeall_unix.go (source SHA-256: `aceeee0c139404a97188802a7f08f93c2f464ce1daf1b0398efe8b327a7d405d`)
- go@go1.26.10/src/os/removeall_windows.go (source SHA-256: `b2d7a7af1f9506e134f4ec87191c2e494931003d2ec17117211990e2ab98a4b0`)
- go@go1.26.10/src/runtime/cgo/clearenv.go (source SHA-256: `0a680faf208c304a76d6a2ecbc2fbb9515e131a1d23ce26050236d72ca9319b4`)
- go@go1.26.10/src/runtime/cgo/windows.go (source SHA-256: `cb1632e83e14cb0b9d08f2dbb50ce653133e4ce96c86f8da589fb4bf50dc4fb3`)
- go@go1.26.10/src/runtime/cgroup_linux.go (source SHA-256: `a2237d59862ebb1fd5f85750663cfe03965463a6d6bf21e5a411df26504f5f7f`)
- go@go1.26.10/src/runtime/cgroup_stubs.go (source SHA-256: `7393e87ba1b216ea9d027286bc3c1dee709d91581b2e0a6540940d22236d19b5`)
- go@go1.26.10/src/runtime/dit.go (source SHA-256: `bed2939b60746bcbf728179c67cdde90b56af3b0b2bbc6dfe7ab052c58a01d4e`)
- go@go1.26.10/src/runtime/fipsbypass.go (source SHA-256: `fb6c3947c6662d94abbd0a328a7ca40f9266dd0d8ce6240d7508abb27e4fa791`)
- go@go1.26.10/src/runtime/hexdump.go (source SHA-256: `0e089a44ff2b7cee6a2d603ee9af805672b69d443ac5b3a58314ac949e0eceaa`)
- go@go1.26.10/src/runtime/linkname_shim.go (source SHA-256: `bc081a57b4608944b79c16c5bb02c5eba258fa6ee5635c427ffd1ec8081891ae`)
- go@go1.26.10/src/runtime/list.go (source SHA-256: `70f32c4b73984ea5a269f73e5e4ffa743cb9ae840b02e3e234f32ea542098972`)
- go@go1.26.10/src/runtime/list_manual.go (source SHA-256: `8ccc132ba4d8e698b50ab7536aed0c1d9a99001bde89bab41223c0e5f0fcef8f`)
- go@go1.26.10/src/runtime/malloc_stubs.go (source SHA-256: `cd762e3218b6465d5b25e49b4473a44dfb8e4632c48be94250cb4ceda7da2075`)
- go@go1.26.10/src/runtime/mgcmark_greenteagc.go (source SHA-256: `34cd31c189bfc011739dd738656f5e698521a736c3b05e2a2f3da49fc6011f3b`)
- go@go1.26.10/src/runtime/os_linux64.go (source SHA-256: `3414c91c0566e953cc5de08191d73fc99a69fa2d156152269032e4351399dd01`)
- go@go1.26.10/src/runtime/preempt_xreg.go (source SHA-256: `23fafbf05a47b495413a70682066a48c8af14e75f51222149c1aa9103eb555a9`)
- go@go1.26.10/src/runtime/runtime_clearenv.go (source SHA-256: `055ee6d0789f8031ffcc1c806e135c678adadcc3e2458fd6fbb844cc7c0f74d1`)
- go@go1.26.10/src/runtime/runtime_noclearenv.go (source SHA-256: `c0e00d499214baf6257577a4abbc506a8429c52e0fde7ad09318b962750f24b4`)
- go@go1.26.10/src/runtime/secret_nosecret.go (source SHA-256: `61ff4c51b846dcd7f2742d43d366ae4ff77af5d68027182be92cfe1e19f826a3`)
- go@go1.26.10/src/runtime/set_vma_name_linux.go (source SHA-256: `d4360e1d3e9ca42ec05150de6820e5ac551b4cb6ad669a4eac50e5da0090659c`)
- go@go1.26.10/src/runtime/set_vma_name_stub.go (source SHA-256: `f06348abd310c60f1e7ee26063bfa46654b129277b2b41bbfae2dad46af0948f`)
- go@go1.26.10/src/runtime/signal_windows_amd64.go (source SHA-256: `bbd73d64a342381fa56a3b3ddc1d7599f88b4128e08c90aaf1881dd5b03c2d73`)
- go@go1.26.10/src/runtime/valgrind0.go (source SHA-256: `a5ba5f53f63ee9b41dffe9003a5c9fffba1530d4397ebe2a4164dbf0e49c1aa4`)
- go@go1.26.10/src/strconv/number.go (source SHA-256: `dda43e8c91d95f9dea0c46b2fec6637e5520cae099e40c5965dca0e41bbaab9e`)
- go@go1.26.10/src/unique/canonmap.go (source SHA-256: `c2231dbb996483905d256f0c1d8cf64266fd89c63bc339e0530e536bedb5aedb`)
- go@go1.26.10/src/vendor/golang.org/x/crypto/chacha20poly1305/fips140only_go1.26.go (source SHA-256: `60368cae6d4d630509309819d7f397532c6f7812978c12983bca375e4ae810aa`)
- go@go1.26.10/src/vendor/golang.org/x/net/dns/dnsmessage/svcb.go (source SHA-256: `2a6051c720262a3ebb63c191f6886dd26b5b049a823345a4c15586395f69b4f2`)

````text
// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/internal/runtime/gc/scan/filter_amd64.s (source SHA-256: `d75bc88f985ffca0e5d1fa0ea53ebc69478b8a4db4bd4618babda1a0f510ab75`)
- go@go1.26.10/src/internal/runtime/gc/scan/scan_amd64.s (source SHA-256: `1123e5160bb6920843be16d6d218a64352b34ce0f6b8ce2df9ed34c68fd642fb`)
- go@go1.26.10/src/internal/runtime/syscall/windows/asm_windows_amd64.s (source SHA-256: `3b99b6a681ce2cd27adc6d4a1b05e841d9cdeb5fff6962936a83c9204d60d5c6`)
- go@go1.26.10/src/math/big/arith_amd64.s (source SHA-256: `388854f0d34f24afcc23c09b83a8dce7a6303c774dc9c0f2295a052ca12f2802`)
- go@go1.26.10/src/math/big/arith_arm64.s (source SHA-256: `1688707cea27c4be6e0ce0555b6bad430dcd20309ae007616c76505d851d5fda`)
- go@go1.26.10/src/runtime/cgo/gcc_clearenv.c (source SHA-256: `9c5e62afce452d425d93793a0dc77921b57bf382b52204f64cefc75cd15b823d`)

````text
// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/drbg/entropy_fips140.go (source SHA-256: `bd3c834a29e31c93e56d81da54d4a874561814a86550de1456c8d6d40a9c5fda`)
- go@go1.26.10/src/internal/poll/fstatat_unix.go (source SHA-256: `77c3abab202397b95ef3b107301bff55f0d18e809982f6af159e6685dbee5a9f`)
- go@go1.26.10/src/internal/syscall/unix/fchmodat_linux.go (source SHA-256: `d3485bf5a7dbf5cf7d9e22bab92730e5973e4ecff48f15f0a4220c924cb3ed73`)
- go@go1.26.10/src/os/exec/lookpath.go (source SHA-256: `8c2b44aacefa3a0b7a92621995cd759396252191e0c5c01b46316df1d948d0d0`)
- go@go1.26.10/src/os/statat.go (source SHA-256: `e10fddfd50fabe4327e2a7eed3b4375ecf006d76dd0319577ddf4b19cdbe99d9`)
- go@go1.26.10/src/os/statat_unix.go (source SHA-256: `0595ce8a8868674b6069431cd182c64b4d24c7b39e38b01f1b108db261d27830`)

````text
// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- go@go1.26.10/src/runtime/memmove_amd64.s (source SHA-256: `1d239c1a965baf43cbdfb82a5144d33aadd95d20ff95becf8347a93ddf08985d`)

````text
// Derived from Inferno's libkern/memmove-386.s (adapted for amd64)
// https://bitbucket.org/inferno-os/inferno-os/src/master/libkern/memmove-386.s
//
//         Copyright © 1994-1999 Lucent Technologies Inc. All rights reserved.
//         Revisions Copyright © 2000-2007 Vita Nuova Holdings Limited (www.vitanuova.com).  All rights reserved.
//         Portions Copyright 2009 The Go Authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/edwards25519/scalar.go (source SHA-256: `b9bc9ba5a2160dccee62ccbb4371c90c5d1c654d55b8ae08f33fa29c58876050`)

````text
// The field implementation in scalar_fiat.go is generated by the fiat-crypto
// project (https://github.com/mit-plv/fiat-crypto) at version v0.0.9 (23d2dbc)
// from a formally verified model.
//
// fiat-crypto code comes under the following license.
//
//     Copyright (c) 2015-2020 The fiat-crypto Authors. All rights reserved.
//
//     Redistribution and use in source and binary forms, with or without
//     modification, are permitted provided that the following conditions are
//     met:
//
//         1. Redistributions of source code must retain the above copyright
//         notice, this list of conditions and the following disclaimer.
//
//     THIS SOFTWARE IS PROVIDED BY the fiat-crypto authors "AS IS"
//     AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO,
//     THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR
//     PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL Berkeley Software Design,
//     Inc. BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
//     EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO,
//     PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR
//     PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
//     LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
//     NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
//     SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
````

### Original source notice

- go@go1.26.10/src/math/cbrt.go (source SHA-256: `de08f8225477bd143bd5318c5f46b7e26bdb008db8e4d174c6130e751e66e335`)

````text
// The go code is a modified version of the original C code from
// http://www.netlib.org/fdlibm/s_cbrt.c and came with this notice.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunSoft, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
````

### Original source notice

- go@go1.26.10/src/math/exp_amd64.s (source SHA-256: `6b0b9ca5b73d27e457fa3ce283dec9fdbc8411ac3b28afe262059330c32e9980`)

````text
// The method is based on a paper by Naoki Shibata: "Efficient evaluation
// methods of elementary functions suitable for SIMD computation", Proc.
// of International Supercomputing Conference 2010 (ISC'10), pp. 25 -- 32
// (May 2010). The paper is available at
// https://link.springer.com/article/10.1007/s00450-010-0108-2
//
// The original code and the constants below are from the author's
// implementation available at http://freshmeat.net/projects/sleef.
// The README file says, "The software is in public domain.
// You can use the software without any obligation."
//
// This code is a simplified version of the original.
````

### Original source notice

- go@go1.26.10/src/math/remainder.go (source SHA-256: `b86740c9a71e4144feec5e728c13f18b157744151d55eeaada90c9c9f26cefef`)

````text
// The original C code and the comment below are from
// FreeBSD's /usr/src/lib/msun/src/e_remainder.c and came
// with this notice. The go code is a simplified version of
// the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_remainder(x,y)
// Return :
//      returns  x REM y  =  x - [x/y]*y  as if in infinite
//      precision arithmetic, where [x/y] is the (infinite bit)
//      integer nearest x/y (in half way cases, choose the even one).
// Method :
//      Based on Mod() returning  x - [x/y]chopped * y  exactly.
````

### Original source notice

- go@go1.26.10/src/math/j0.go (source SHA-256: `8d226442c38eeb79e5273980fc5d94f0ba9d4c4407d0e9a7e8db45838fc128f8`)

````text
// The original C code and the long comment below are
// from FreeBSD's /usr/src/lib/msun/src/e_j0.c and
// came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_j0(x), __ieee754_y0(x)
// Bessel function of the first and second kinds of order zero.
// Method -- j0(x):
//      1. For tiny x, we use j0(x) = 1 - x**2/4 + x**4/64 - ...
//      2. Reduce x to |x| since j0(x)=j0(-x),  and
//         for x in (0,2)
//              j0(x) = 1-z/4+ z**2*R0/S0,  where z = x*x;
//         (precision:  |j0-1+z/4-z**2R0/S0 |<2**-63.67 )
//         for x in (2,inf)
//              j0(x) = sqrt(2/(pi*x))*(p0(x)*cos(x0)-q0(x)*sin(x0))
//         where x0 = x-pi/4. It is better to compute sin(x0),cos(x0)
//         as follow:
//              cos(x0) = cos(x)cos(pi/4)+sin(x)sin(pi/4)
//                      = 1/sqrt(2) * (cos(x) + sin(x))
//              sin(x0) = sin(x)cos(pi/4)-cos(x)sin(pi/4)
//                      = 1/sqrt(2) * (sin(x) - cos(x))
//         (To avoid cancellation, use
//              sin(x) +- cos(x) = -cos(2x)/(sin(x) -+ cos(x))
//         to compute the worse one.)
//
//      3 Special cases
//              j0(nan)= nan
//              j0(0) = 1
//              j0(inf) = 0
//
// Method -- y0(x):
//      1. For x<2.
//         Since
//              y0(x) = 2/pi*(j0(x)*(ln(x/2)+Euler) + x**2/4 - ...)
//         therefore y0(x)-2/pi*j0(x)*ln(x) is an even function.
//         We use the following function to approximate y0,
//              y0(x) = U(z)/V(z) + (2/pi)*(j0(x)*ln(x)), z= x**2
//         where
//              U(z) = u00 + u01*z + ... + u06*z**6
//              V(z) = 1  + v01*z + ... + v04*z**4
//         with absolute approximation error bounded by 2**-72.
//         Note: For tiny x, U/V = u0 and j0(x)~1, hence
//              y0(tiny) = u0 + (2/pi)*ln(tiny), (choose tiny<2**-27)
//      2. For x>=2.
//              y0(x) = sqrt(2/(pi*x))*(p0(x)*cos(x0)+q0(x)*sin(x0))
//         where x0 = x-pi/4. It is better to compute sin(x0),cos(x0)
//         by the method mentioned above.
//      3. Special cases: y0(0)=-inf, y0(x<0)=NaN, y0(inf)=0.
//
````

### Original source notice

- go@go1.26.10/src/math/j1.go (source SHA-256: `599ddf54522dd0cf19e3cdcee1298a5ca98eaa529e6e647fb8311373f5bef1fa`)

````text
// The original C code and the long comment below are
// from FreeBSD's /usr/src/lib/msun/src/e_j1.c and
// came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_j1(x), __ieee754_y1(x)
// Bessel function of the first and second kinds of order one.
// Method -- j1(x):
//      1. For tiny x, we use j1(x) = x/2 - x**3/16 + x**5/384 - ...
//      2. Reduce x to |x| since j1(x)=-j1(-x),  and
//         for x in (0,2)
//              j1(x) = x/2 + x*z*R0/S0,  where z = x*x;
//         (precision:  |j1/x - 1/2 - R0/S0 |<2**-61.51 )
//         for x in (2,inf)
//              j1(x) = sqrt(2/(pi*x))*(p1(x)*cos(x1)-q1(x)*sin(x1))
//              y1(x) = sqrt(2/(pi*x))*(p1(x)*sin(x1)+q1(x)*cos(x1))
//         where x1 = x-3*pi/4. It is better to compute sin(x1),cos(x1)
//         as follow:
//              cos(x1) =  cos(x)cos(3pi/4)+sin(x)sin(3pi/4)
//                      =  1/sqrt(2) * (sin(x) - cos(x))
//              sin(x1) =  sin(x)cos(3pi/4)-cos(x)sin(3pi/4)
//                      = -1/sqrt(2) * (sin(x) + cos(x))
//         (To avoid cancellation, use
//              sin(x) +- cos(x) = -cos(2x)/(sin(x) -+ cos(x))
//         to compute the worse one.)
//
//      3 Special cases
//              j1(nan)= nan
//              j1(0) = 0
//              j1(inf) = 0
//
// Method -- y1(x):
//      1. screen out x<=0 cases: y1(0)=-inf, y1(x<0)=NaN
//      2. For x<2.
//         Since
//              y1(x) = 2/pi*(j1(x)*(ln(x/2)+Euler)-1/x-x/2+5/64*x**3-...)
//         therefore y1(x)-2/pi*j1(x)*ln(x)-1/x is an odd function.
//         We use the following function to approximate y1,
//              y1(x) = x*U(z)/V(z) + (2/pi)*(j1(x)*ln(x)-1/x), z= x**2
//         where for x in [0,2] (abs err less than 2**-65.89)
//              U(z) = U0[0] + U0[1]*z + ... + U0[4]*z**4
//              V(z) = 1  + v0[0]*z + ... + v0[4]*z**5
//         Note: For tiny x, 1/x dominate y1 and hence
//              y1(tiny) = -2/pi/tiny, (choose tiny<2**-54)
//      3. For x>=2.
//               y1(x) = sqrt(2/(pi*x))*(p1(x)*sin(x1)+q1(x)*cos(x1))
//         where x1 = x-3*pi/4. It is better to compute sin(x1),cos(x1)
//         by method mentioned above.
````

### Original source notice

- go@go1.26.10/src/math/jn.go (source SHA-256: `b588c5389aa3ac04eeb72b29a610bef4af02b8ee7faf3d50717d2da2219bc144`)

````text
// The original C code and the long comment below are
// from FreeBSD's /usr/src/lib/msun/src/e_jn.c and
// came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_jn(n, x), __ieee754_yn(n, x)
// floating point Bessel's function of the 1st and 2nd kind
// of order n
//
// Special cases:
//      y0(0)=y1(0)=yn(n,0) = -inf with division by zero signal;
//      y0(-ve)=y1(-ve)=yn(n,-ve) are NaN with invalid signal.
// Note 2. About jn(n,x), yn(n,x)
//      For n=0, j0(x) is called,
//      for n=1, j1(x) is called,
//      for n<x, forward recursion is used starting
//      from values of j0(x) and j1(x).
//      for n>x, a continued fraction approximation to
//      j(n,x)/j(n-1,x) is evaluated and then backward
//      recursion is used starting from a supposed value
//      for j(n,x). The resulting value of j(0,x) is
//      compared with the actual value to correct the
//      supposed value of j(n,x).
//
//      yn(n,x) is similar in all respects, except
//      that forward recursion is used for all
//      values of n>1.
````

### Original source notice

- go@go1.26.10/src/math/lgamma.go (source SHA-256: `c07c6b7b6d57e97f47531abca86accbd23e0a0d49d8d5cb60c226e87df9c6804`)

````text
// The original C code and the long comment below are
// from FreeBSD's /usr/src/lib/msun/src/e_lgamma_r.c and
// came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_lgamma_r(x, signgamp)
// Reentrant version of the logarithm of the Gamma function
// with user provided pointer for the sign of Gamma(x).
//
// Method:
//   1. Argument Reduction for 0 < x <= 8
//      Since gamma(1+s)=s*gamma(s), for x in [0,8], we may
//      reduce x to a number in [1.5,2.5] by
//              lgamma(1+s) = log(s) + lgamma(s)
//      for example,
//              lgamma(7.3) = log(6.3) + lgamma(6.3)
//                          = log(6.3*5.3) + lgamma(5.3)
//                          = log(6.3*5.3*4.3*3.3*2.3) + lgamma(2.3)
//   2. Polynomial approximation of lgamma around its
//      minimum (ymin=1.461632144968362245) to maintain monotonicity.
//      On [ymin-0.23, ymin+0.27] (i.e., [1.23164,1.73163]), use
//              Let z = x-ymin;
//              lgamma(x) = -1.214862905358496078218 + z**2*poly(z)
//              poly(z) is a 14 degree polynomial.
//   2. Rational approximation in the primary interval [2,3]
//      We use the following approximation:
//              s = x-2.0;
//              lgamma(x) = 0.5*s + s*P(s)/Q(s)
//      with accuracy
//              |P/Q - (lgamma(x)-0.5s)| < 2**-61.71
//      Our algorithms are based on the following observation
//
//                             zeta(2)-1    2    zeta(3)-1    3
// lgamma(2+s) = s*(1-Euler) + --------- * s  -  --------- * s  + ...
//                                 2                 3
//
//      where Euler = 0.5772156649... is the Euler constant, which
//      is very close to 0.5.
//
//   3. For x>=8, we have
//      lgamma(x)~(x-0.5)log(x)-x+0.5*log(2pi)+1/(12x)-1/(360x**3)+....
//      (better formula:
//         lgamma(x)~(x-0.5)*(log(x)-1)-.5*(log(2pi)-1) + ...)
//      Let z = 1/x, then we approximation
//              f(z) = lgamma(x) - (x-0.5)(log(x)-1)
//      by
//                                  3       5             11
//              w = w0 + w1*z + w2*z  + w3*z  + ... + w6*z
//      where
//              |w - f(z)| < 2**-58.74
//
//   4. For negative x, since (G is gamma function)
//              -x*G(-x)*G(x) = pi/sin(pi*x),
//      we have
//              G(x) = pi/(sin(pi*x)*(-x)*G(-x))
//      since G(-x) is positive, sign(G(x)) = sign(sin(pi*x)) for x<0
//      Hence, for x<0, signgam = sign(sin(pi*x)) and
//              lgamma(x) = log(|Gamma(x)|)
//                        = log(pi/(|x*sin(pi*x)|)) - lgamma(-x);
//      Note: one should avoid computing pi*(-x) directly in the
//            computation of sin(pi*(-x)).
//
//   5. Special Cases
//              lgamma(2+s) ~ s*(1-Euler) for tiny s
//              lgamma(1)=lgamma(2)=0
//              lgamma(x) ~ -log(x) for tiny x
//              lgamma(0) = lgamma(inf) = inf
//              lgamma(-integer) = +-inf
//
//
````

### Original source notice

- go@go1.26.10/src/math/sqrt.go (source SHA-256: `c7e036d1511bea94992014f9b8c9259e5737c385cdd5135db9092369250d4a94`)

````text
// The original C code and the long comment below are
// from FreeBSD's /usr/src/lib/msun/src/e_sqrt.c and
// came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_sqrt(x)
// Return correctly rounded sqrt.
//           -----------------------------------------
//           | Use the hardware sqrt if you have one |
//           -----------------------------------------
// Method:
//   Bit by bit method using integer arithmetic. (Slow, but portable)
//   1. Normalization
//      Scale x to y in [1,4) with even powers of 2:
//      find an integer k such that  1 <= (y=x*2**(2k)) < 4, then
//              sqrt(x) = 2**k * sqrt(y)
//   2. Bit by bit computation
//      Let q  = sqrt(y) truncated to i bit after binary point (q = 1),
//           i                                                   0
//                                     i+1         2
//          s  = 2*q , and      y  =  2   * ( y - q  ).          (1)
//           i      i            i                 i
//
//      To compute q    from q , one checks whether
//                  i+1       i
//
//                            -(i+1) 2
//                      (q + 2      )  <= y.                     (2)
//                        i
//                                                            -(i+1)
//      If (2) is false, then q   = q ; otherwise q   = q  + 2      .
//                             i+1   i             i+1   i
//
//      With some algebraic manipulation, it is not difficult to see
//      that (2) is equivalent to
//                             -(i+1)
//                      s  +  2       <= y                       (3)
//                       i                i
//
//      The advantage of (3) is that s  and y  can be computed by
//                                    i      i
//      the following recurrence formula:
//          if (3) is false
//
//          s     =  s  ,       y    = y   ;                     (4)
//           i+1      i          i+1    i
//
//      otherwise,
//                         -i                      -(i+1)
//          s     =  s  + 2  ,  y    = y  -  s  - 2              (5)
//           i+1      i          i+1    i     i
//
//      One may easily use induction to prove (4) and (5).
//      Note. Since the left hand side of (3) contain only i+2 bits,
//            it is not necessary to do a full (53-bit) comparison
//            in (3).
//   3. Final rounding
//      After generating the 53 bits result, we compute one more bit.
//      Together with the remainder, we can decide whether the
//      result is exact, bigger than 1/2ulp, or less than 1/2ulp
//      (it will never equal to 1/2ulp).
//      The rounding mode can be detected by checking whether
//      huge + tiny is equal to huge, and whether huge - tiny is
//      equal to huge for some floating point number "huge" and "tiny".
//
//
// Notes:  Rounding mode detection omitted. The constants "mask", "shift",
// and "bias" are found in src/math/bits.go
````

### Original source notice

- go@go1.26.10/src/math/erf.go (source SHA-256: `31cba8a8f66d6081eb393b85e6a0d2a3737ebb3fc3fb192867e393a2907c0c8c`)

````text
// The original C code and the long comment below are
// from FreeBSD's /usr/src/lib/msun/src/s_erf.c and
// came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// double erf(double x)
// double erfc(double x)
//                           x
//                    2      |\
//     erf(x)  =  ---------  | exp(-t*t)dt
//                 sqrt(pi) \|
//                           0
//
//     erfc(x) =  1-erf(x)
//  Note that
//              erf(-x) = -erf(x)
//              erfc(-x) = 2 - erfc(x)
//
// Method:
//      1. For |x| in [0, 0.84375]
//          erf(x)  = x + x*R(x**2)
//          erfc(x) = 1 - erf(x)           if x in [-.84375,0.25]
//                  = 0.5 + ((0.5-x)-x*R)  if x in [0.25,0.84375]
//         where R = P/Q where P is an odd poly of degree 8 and
//         Q is an odd poly of degree 10.
//                                               -57.90
//                      | R - (erf(x)-x)/x | <= 2
//
//
//         Remark. The formula is derived by noting
//          erf(x) = (2/sqrt(pi))*(x - x**3/3 + x**5/10 - x**7/42 + ....)
//         and that
//          2/sqrt(pi) = 1.128379167095512573896158903121545171688
//         is close to one. The interval is chosen because the fix
//         point of erf(x) is near 0.6174 (i.e., erf(x)=x when x is
//         near 0.6174), and by some experiment, 0.84375 is chosen to
//         guarantee the error is less than one ulp for erf.
//
//      2. For |x| in [0.84375,1.25], let s = |x| - 1, and
//         c = 0.84506291151 rounded to single (24 bits)
//              erf(x)  = sign(x) * (c  + P1(s)/Q1(s))
//              erfc(x) = (1-c)  - P1(s)/Q1(s) if x > 0
//                        1+(c+P1(s)/Q1(s))    if x < 0
//              |P1/Q1 - (erf(|x|)-c)| <= 2**-59.06
//         Remark: here we use the taylor series expansion at x=1.
//              erf(1+s) = erf(1) + s*Poly(s)
//                       = 0.845.. + P1(s)/Q1(s)
//         That is, we use rational approximation to approximate
//                      erf(1+s) - (c = (single)0.84506291151)
//         Note that |P1/Q1|< 0.078 for x in [0.84375,1.25]
//         where
//              P1(s) = degree 6 poly in s
//              Q1(s) = degree 6 poly in s
//
//      3. For x in [1.25,1/0.35(~2.857143)],
//              erfc(x) = (1/x)*exp(-x*x-0.5625+R1/S1)
//              erf(x)  = 1 - erfc(x)
//         where
//              R1(z) = degree 7 poly in z, (z=1/x**2)
//              S1(z) = degree 8 poly in z
//
//      4. For x in [1/0.35,28]
//              erfc(x) = (1/x)*exp(-x*x-0.5625+R2/S2) if x > 0
//                      = 2.0 - (1/x)*exp(-x*x-0.5625+R2/S2) if -6<x<0
//                      = 2.0 - tiny            (if x <= -6)
//              erf(x)  = sign(x)*(1.0 - erfc(x)) if x < 6, else
//              erf(x)  = sign(x)*(1.0 - tiny)
//         where
//              R2(z) = degree 6 poly in z, (z=1/x**2)
//              S2(z) = degree 7 poly in z
//
//      Note1:
//         To compute exp(-x*x-0.5625+R/S), let s be a single
//         precision number and s := x; then
//              -x*x = -s*s + (s-x)*(s+x)
//              exp(-x*x-0.5626+R/S) =
//                      exp(-s*s-0.5625)*exp((s-x)*(s+x)+R/S);
//      Note2:
//         Here 4 and 5 make use of the asymptotic series
//                        exp(-x*x)
//              erfc(x) ~ ---------- * ( 1 + Poly(1/x**2) )
//                        x*sqrt(pi)
//         We use rational approximation to approximate
//              g(s)=f(1/x**2) = log(erfc(x)*x) - x*x + 0.5625
//         Here is the error bound for R1/S1 and R2/S2
//              |R1/S1 - f(x)|  < 2**(-62.57)
//              |R2/S2 - f(x)|  < 2**(-61.52)
//
//      5. For inf > x >= 28
//              erf(x)  = sign(x) *(1 - tiny)  (raise inexact)
//              erfc(x) = tiny*tiny (raise underflow) if x > 0
//                      = 2 - tiny if x<0
//
//      7. Special case:
//              erf(0)  = 0, erf(inf)  = 1, erf(-inf) = -1,
//              erfc(0) = 1, erfc(inf) = 0, erfc(-inf) = 2,
//              erfc/erf(NaN) is NaN
````

### Original source notice

- go@go1.26.10/src/math/acosh.go (source SHA-256: `e1b20885bef3ccb7763ae868d23120194cf38faef70f1af45d2870ccb1a00140`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_acosh.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// __ieee754_acosh(x)
// Method :
//	Based on
//	        acosh(x) = log [ x + sqrt(x*x-1) ]
//	we have
//	        acosh(x) := log(x)+ln2,	if x is large; else
//	        acosh(x) := log(2x-1/(sqrt(x*x-1)+x)) if x>2; else
//	        acosh(x) := log1p(t+sqrt(2.0*t+t*t)); where t=x-1.
//
// Special cases:
//	acosh(x) is NaN with signal if x<1.
//	acosh(NaN) is NaN without signal.
//
````

### Original source notice

- go@go1.26.10/src/math/atanh.go (source SHA-256: `2afbcd791957596e7af481c3d23640f15e2d1149d8ed70fc108ddd1d5ab79430`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_atanh.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// __ieee754_atanh(x)
// Method :
//	1. Reduce x to positive by atanh(-x) = -atanh(x)
//	2. For x>=0.5
//	            1              2x                          x
//	atanh(x) = --- * log(1 + -------) = 0.5 * log1p(2 * --------)
//	            2             1 - x                      1 - x
//
//	For x<0.5
//	atanh(x) = 0.5*log1p(2x+2x*x/(1-x))
//
// Special cases:
//	atanh(x) is NaN if |x| > 1 with signal;
//	atanh(NaN) is that NaN with no signal;
//	atanh(+-1) is +-INF with signal.
//
````

### Original source notice

- go@go1.26.10/src/math/exp.go (source SHA-256: `75d415ac8eca1afc796cd621b393153ff82d342c24f5b391dfdd6a6e70303f2b`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_exp.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 2004 by Sun Microsystems, Inc. All rights reserved.
//
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// exp(x)
// Returns the exponential of x.
//
// Method
//   1. Argument reduction:
//      Reduce x to an r so that |r| <= 0.5*ln2 ~ 0.34658.
//      Given x, find r and integer k such that
//
//               x = k*ln2 + r,  |r| <= 0.5*ln2.
//
//      Here r will be represented as r = hi-lo for better
//      accuracy.
//
//   2. Approximation of exp(r) by a special rational function on
//      the interval [0,0.34658]:
//      Write
//          R(r**2) = r*(exp(r)+1)/(exp(r)-1) = 2 + r*r/6 - r**4/360 + ...
//      We use a special Remez algorithm on [0,0.34658] to generate
//      a polynomial of degree 5 to approximate R. The maximum error
//      of this polynomial approximation is bounded by 2**-59. In
//      other words,
//          R(z) ~ 2.0 + P1*z + P2*z**2 + P3*z**3 + P4*z**4 + P5*z**5
//      (where z=r*r, and the values of P1 to P5 are listed below)
//      and
//          |                  5          |     -59
//          | 2.0+P1*z+...+P5*z   -  R(z) | <= 2
//          |                             |
//      The computation of exp(r) thus becomes
//                             2*r
//              exp(r) = 1 + -------
//                            R - r
//                                 r*R1(r)
//                     = 1 + r + ----------- (for better accuracy)
//                                2 - R1(r)
//      where
//                               2       4             10
//              R1(r) = r - (P1*r  + P2*r  + ... + P5*r   ).
//
//   3. Scale back to obtain exp(x):
//      From step 1, we have
//         exp(x) = 2**k * exp(r)
//
// Special cases:
//      exp(INF) is INF, exp(NaN) is NaN;
//      exp(-INF) is 0, and
//      for finite argument, only exp(0)=1 is exact.
//
// Accuracy:
//      according to an error analysis, the error is always less than
//      1 ulp (unit in the last place).
//
// Misc. info.
//      For IEEE double
//          if x >  7.09782712893383973096e+02 then exp(x) overflow
//          if x < -7.45133219101941108420e+02 then exp(x) underflow
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.
````

### Original source notice

- go@go1.26.10/src/math/log.go (source SHA-256: `d8675d3d9c665ec0db5cc15478c03efb244301bfa1268d02ae1861bf6f2812ab`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/e_log.c
// and came with this notice. The go code is a simpler
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// __ieee754_log(x)
// Return the logarithm of x
//
// Method :
//   1. Argument Reduction: find k and f such that
//			x = 2**k * (1+f),
//	   where  sqrt(2)/2 < 1+f < sqrt(2) .
//
//   2. Approximation of log(1+f).
//	Let s = f/(2+f) ; based on log(1+f) = log(1+s) - log(1-s)
//		 = 2s + 2/3 s**3 + 2/5 s**5 + .....,
//	     	 = 2s + s*R
//      We use a special Reme algorithm on [0,0.1716] to generate
//	a polynomial of degree 14 to approximate R.  The maximum error
//	of this polynomial approximation is bounded by 2**-58.45. In
//	other words,
//		        2      4      6      8      10      12      14
//	    R(z) ~ L1*s +L2*s +L3*s +L4*s +L5*s  +L6*s  +L7*s
//	(the values of L1 to L7 are listed in the program) and
//	    |      2          14          |     -58.45
//	    | L1*s +...+L7*s    -  R(z) | <= 2
//	    |                             |
//	Note that 2s = f - s*f = f - hfsq + s*hfsq, where hfsq = f*f/2.
//	In order to guarantee error in log below 1ulp, we compute log by
//		log(1+f) = f - s*(f - R)		(if f is not too large)
//		log(1+f) = f - (hfsq - s*(hfsq+R)).	(better accuracy)
//
//	3. Finally,  log(x) = k*Ln2 + log(1+f).
//			    = k*Ln2_hi+(f-(hfsq-(s*(hfsq+R)+k*Ln2_lo)))
//	   Here Ln2 is split into two floating point number:
//			Ln2_hi + Ln2_lo,
//	   where n*Ln2_hi is always exact for |n| < 2000.
//
// Special cases:
//	log(x) is NaN with signal if x < 0 (including -INF) ;
//	log(+INF) is +INF; log(0) is -INF with signal;
//	log(NaN) is that NaN with no signal.
//
// Accuracy:
//	according to an error analysis, the error is always less than
//	1 ulp (unit in the last place).
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.
````

### Original source notice

- go@go1.26.10/src/math/asinh.go (source SHA-256: `7bd3753d8bb94eb86682c9215c6266efccd75a66701aad16fc94d9fbfb89ccb1`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/s_asinh.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// asinh(x)
// Method :
//	Based on
//	        asinh(x) = sign(x) * log [ |x| + sqrt(x*x+1) ]
//	we have
//	asinh(x) := x  if  1+x*x=1,
//	         := sign(x)*(log(x)+ln2) for large |x|, else
//	         := sign(x)*log(2|x|+1/(|x|+sqrt(x*x+1))) if|x|>2, else
//	         := sign(x)*log1p(|x| + x**2/(1 + sqrt(1+x**2)))
//
````

### Original source notice

- go@go1.26.10/src/math/expm1.go (source SHA-256: `b05c184f92e2e4bdfb09eb34e6d2474fe4570a4eaf5eb1e128d2b6904908d63c`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/s_expm1.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
// expm1(x)
// Returns exp(x)-1, the exponential of x minus 1.
//
// Method
//   1. Argument reduction:
//      Given x, find r and integer k such that
//
//               x = k*ln2 + r,  |r| <= 0.5*ln2 ~ 0.34658
//
//      Here a correction term c will be computed to compensate
//      the error in r when rounded to a floating-point number.
//
//   2. Approximating expm1(r) by a special rational function on
//      the interval [0,0.34658]:
//      Since
//          r*(exp(r)+1)/(exp(r)-1) = 2+ r**2/6 - r**4/360 + ...
//      we define R1(r*r) by
//          r*(exp(r)+1)/(exp(r)-1) = 2+ r**2/6 * R1(r*r)
//      That is,
//          R1(r**2) = 6/r *((exp(r)+1)/(exp(r)-1) - 2/r)
//                   = 6/r * ( 1 + 2.0*(1/(exp(r)-1) - 1/r))
//                   = 1 - r**2/60 + r**4/2520 - r**6/100800 + ...
//      We use a special Reme algorithm on [0,0.347] to generate
//      a polynomial of degree 5 in r*r to approximate R1. The
//      maximum error of this polynomial approximation is bounded
//      by 2**-61. In other words,
//          R1(z) ~ 1.0 + Q1*z + Q2*z**2 + Q3*z**3 + Q4*z**4 + Q5*z**5
//      where   Q1  =  -1.6666666666666567384E-2,
//              Q2  =   3.9682539681370365873E-4,
//              Q3  =  -9.9206344733435987357E-6,
//              Q4  =   2.5051361420808517002E-7,
//              Q5  =  -6.2843505682382617102E-9;
//      (where z=r*r, and the values of Q1 to Q5 are listed below)
//      with error bounded by
//          |                  5           |     -61
//          | 1.0+Q1*z+...+Q5*z   -  R1(z) | <= 2
//          |                              |
//
//      expm1(r) = exp(r)-1 is then computed by the following
//      specific way which minimize the accumulation rounding error:
//                             2     3
//                            r     r    [ 3 - (R1 + R1*r/2)  ]
//            expm1(r) = r + --- + --- * [--------------------]
//                            2     2    [ 6 - r*(3 - R1*r/2) ]
//
//      To compensate the error in the argument reduction, we use
//              expm1(r+c) = expm1(r) + c + expm1(r)*c
//                         ~ expm1(r) + c + r*c
//      Thus c+r*c will be added in as the correction terms for
//      expm1(r+c). Now rearrange the term to avoid optimization
//      screw up:
//                      (      2                                    2 )
//                      ({  ( r    [ R1 -  (3 - R1*r/2) ]  )  }    r  )
//       expm1(r+c)~r - ({r*(--- * [--------------------]-c)-c} - --- )
//                      ({  ( 2    [ 6 - r*(3 - R1*r/2) ]  )  }    2  )
//                      (                                             )
//
//                 = r - E
//   3. Scale back to obtain expm1(x):
//      From step 1, we have
//         expm1(x) = either 2**k*[expm1(r)+1] - 1
//                  = or     2**k*[expm1(r) + (1-2**-k)]
//   4. Implementation notes:
//      (A). To save one multiplication, we scale the coefficient Qi
//           to Qi*2**i, and replace z by (x**2)/2.
//      (B). To achieve maximum accuracy, we compute expm1(x) by
//        (i)   if x < -56*ln2, return -1.0, (raise inexact if x!=inf)
//        (ii)  if k=0, return r-E
//        (iii) if k=-1, return 0.5*(r-E)-0.5
//        (iv)  if k=1 if r < -0.25, return 2*((r+0.5)- E)
//                     else          return  1.0+2.0*(r-E);
//        (v)   if (k<-2||k>56) return 2**k(1-(E-r)) - 1 (or exp(x)-1)
//        (vi)  if k <= 20, return 2**k((1-2**-k)-(E-r)), else
//        (vii) return 2**k(1-((E+2**-k)-r))
//
// Special cases:
//      expm1(INF) is INF, expm1(NaN) is NaN;
//      expm1(-INF) is -1, and
//      for finite argument, only expm1(0)=0 is exact.
//
// Accuracy:
//      according to an error analysis, the error is always less than
//      1 ulp (unit in the last place).
//
// Misc. info.
//      For IEEE double
//          if x >  7.09782712893383973096e+02 then expm1(x) overflow
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.
//
````

### Original source notice

- go@go1.26.10/src/math/log1p.go (source SHA-256: `10cb153951ada1c8c5ed3897caded00d90a7314c0640812847c3064a0c5f4b06`)

````text
// The original C code, the long comment, and the constants
// below are from FreeBSD's /usr/src/lib/msun/src/s_log1p.c
// and came with this notice. The go code is a simplified
// version of the original C.
//
// ====================================================
// Copyright (C) 1993 by Sun Microsystems, Inc. All rights reserved.
//
// Developed at SunPro, a Sun Microsystems, Inc. business.
// Permission to use, copy, modify, and distribute this
// software is freely granted, provided that this notice
// is preserved.
// ====================================================
//
//
// double log1p(double x)
//
// Method :
//   1. Argument Reduction: find k and f such that
//                      1+x = 2**k * (1+f),
//         where  sqrt(2)/2 < 1+f < sqrt(2) .
//
//      Note. If k=0, then f=x is exact. However, if k!=0, then f
//      may not be representable exactly. In that case, a correction
//      term is need. Let u=1+x rounded. Let c = (1+x)-u, then
//      log(1+x) - log(u) ~ c/u. Thus, we proceed to compute log(u),
//      and add back the correction term c/u.
//      (Note: when x > 2**53, one can simply return log(x))
//
//   2. Approximation of log1p(f).
//      Let s = f/(2+f) ; based on log(1+f) = log(1+s) - log(1-s)
//               = 2s + 2/3 s**3 + 2/5 s**5 + .....,
//               = 2s + s*R
//      We use a special Reme algorithm on [0,0.1716] to generate
//      a polynomial of degree 14 to approximate R The maximum error
//      of this polynomial approximation is bounded by 2**-58.45. In
//      other words,
//                      2      4      6      8      10      12      14
//          R(z) ~ Lp1*s +Lp2*s +Lp3*s +Lp4*s +Lp5*s  +Lp6*s  +Lp7*s
//      (the values of Lp1 to Lp7 are listed in the program)
//      and
//          |      2          14          |     -58.45
//          | Lp1*s +...+Lp7*s    -  R(z) | <= 2
//          |                             |
//      Note that 2s = f - s*f = f - hfsq + s*hfsq, where hfsq = f*f/2.
//      In order to guarantee error in log below 1ulp, we compute log
//      by
//              log1p(f) = f - (hfsq - s*(hfsq+R)).
//
//   3. Finally, log1p(x) = k*ln2 + log1p(f).
//                        = k*ln2_hi+(f-(hfsq-(s*(hfsq+R)+k*ln2_lo)))
//      Here ln2 is split into two floating point number:
//                   ln2_hi + ln2_lo,
//      where n*ln2_hi is always exact for |n| < 2000.
//
// Special cases:
//      log1p(x) is NaN with signal if x < -1 (including -INF) ;
//      log1p(+INF) is +INF; log1p(-1) is -INF with signal;
//      log1p(NaN) is that NaN with no signal.
//
// Accuracy:
//      according to an error analysis, the error is always less than
//      1 ulp (unit in the last place).
//
// Constants:
// The hexadecimal values are the intended ones for the following
// constants. The decimal values may be used, provided that the
// compiler will convert from decimal to binary accurately enough
// to produce the hexadecimal values shown.
//
// Note: Assuming log() return accurate answer, the following
//       algorithm can be used to compute log1p(x) to within a few ULP:
//
//              u = 1+x;
//              if(u==1.0) return x ; else
//                         return log(u)*(x/(u-1.0));
//
//       See HP-15C Advanced Functions Handbook, p.193.
````

### Original source notice

- go@go1.26.10/src/math/cmplx/asin.go (source SHA-256: `10d60a9dc9aef2e0217e73bf648ae594c91e194de35dfb2eec96aff6e2a38b18`)
- go@go1.26.10/src/math/cmplx/exp.go (source SHA-256: `1660975b05851b22747bffa9e1e8990c73f86a558e364c46b6b3054fdfc596a4`)
- go@go1.26.10/src/math/cmplx/log.go (source SHA-256: `4e78d54adda6eacf3ae1e5986e428b1b946d158657ef781b218694fb578cc85e`)
- go@go1.26.10/src/math/cmplx/pow.go (source SHA-256: `a64cf1d1955098bfec0b2da0cf0122ca5777e33c9f623fcc03f49d950e90e3ca`)
- go@go1.26.10/src/math/cmplx/sin.go (source SHA-256: `de24e3dd23fea2fd71d5206136f4c7786568d785790af1cfae01ef0002c821bd`)
- go@go1.26.10/src/math/cmplx/sqrt.go (source SHA-256: `d55bfa6adb2448163ca500b9f36175e18915c6d9973a7fb117e6c6acc76aecc5`)
- go@go1.26.10/src/math/cmplx/tan.go (source SHA-256: `0afd68f9317f3a278936d7ac50fe0614a6bee7d1a7a0023ffcc7caa25ba1c660`)

````text
// The original C code, the long comment, and the constants
// below are from http://netlib.sandia.gov/cephes/c9x-complex/clog.c.
// The go code is a simplified version of the original C.
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov
````

### Original source notice

- go@go1.26.10/src/math/gamma.go (source SHA-256: `627e52bb7a984ab2da61da481729f317eae7a320a2980c72aa6077ca9e6bd03b`)

````text
// The original C code, the long comment, and the constants
// below are from http://netlib.sandia.gov/cephes/cprob/gamma.c.
// The go code is a simplified version of the original C.
//
//      tgamma.c
//
//      Gamma function
//
// SYNOPSIS:
//
// double x, y, tgamma();
// extern int signgam;
//
// y = tgamma( x );
//
// DESCRIPTION:
//
// Returns gamma function of the argument. The result is
// correctly signed, and the sign (+1 or -1) is also
// returned in a global (extern) variable named signgam.
// This variable is also filled in by the logarithmic gamma
// function lgamma().
//
// Arguments |x| <= 34 are reduced by recurrence and the function
// approximated by a rational function of degree 6/7 in the
// interval (2,3).  Large arguments are handled by Stirling's
// formula. Large negative arguments are made positive using
// a reflection formula.
//
// ACCURACY:
//
//                      Relative error:
// arithmetic   domain     # trials      peak         rms
//    DEC      -34, 34      10000       1.3e-16     2.5e-17
//    IEEE    -170,-33      20000       2.3e-15     3.3e-16
//    IEEE     -33,  33     20000       9.4e-16     2.2e-16
//    IEEE      33, 171.6   20000       2.3e-15     3.2e-16
//
// Error for arguments outside the test range will be larger
// owing to error amplification by the exponential function.
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov
````

### Original source notice

- go@go1.26.10/src/math/sin.go (source SHA-256: `913640c10f5093e8cdaafa1ebebbf5175b5a156bdc0fa89d77350dd4dad3be9b`)

````text
// The original C code, the long comment, and the constants
// below were from http://netlib.sandia.gov/cephes/cmath/sin.c,
// available from http://www.netlib.org/cephes/cmath.tgz.
// The go code is a simplified version of the original C.
//
//      sin.c
//
//      Circular sine
//
// SYNOPSIS:
//
// double x, y, sin();
// y = sin( x );
//
// DESCRIPTION:
//
// Range reduction is into intervals of pi/4.  The reduction error is nearly
// eliminated by contriving an extended precision modular arithmetic.
//
// Two polynomial approximating functions are employed.
// Between 0 and pi/4 the sine is approximated by
//      x  +  x**3 P(x**2).
// Between pi/4 and pi/2 the cosine is represented as
//      1  -  x**2 Q(x**2).
//
// ACCURACY:
//
//                      Relative error:
// arithmetic   domain      # trials      peak         rms
//    DEC       0, 10       150000       3.0e-17     7.8e-18
//    IEEE -1.07e9,+1.07e9  130000       2.1e-16     5.4e-17
//
// Partial loss of accuracy begins to occur at x = 2**30 = 1.074e9.  The loss
// is not gradual, but jumps suddenly to about 1 part in 10e7.  Results may
// be meaningless for x > 2**49 = 5.6e14.
//
//      cos.c
//
//      Circular cosine
//
// SYNOPSIS:
//
// double x, y, cos();
// y = cos( x );
//
// DESCRIPTION:
//
// Range reduction is into intervals of pi/4.  The reduction error is nearly
// eliminated by contriving an extended precision modular arithmetic.
//
// Two polynomial approximating functions are employed.
// Between 0 and pi/4 the cosine is approximated by
//      1  -  x**2 Q(x**2).
// Between pi/4 and pi/2 the sine is represented as
//      x  +  x**3 P(x**2).
//
// ACCURACY:
//
//                      Relative error:
// arithmetic   domain      # trials      peak         rms
//    IEEE -1.07e9,+1.07e9  130000       2.1e-16     5.4e-17
//    DEC        0,+1.07e9   17000       3.0e-17     7.2e-18
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov
````

### Original source notice

- go@go1.26.10/src/math/tan.go (source SHA-256: `b1a549f6d267acf37399b29827ecc4cc152cfb655d8a2ed6dfe56a67e58f8e60`)

````text
// The original C code, the long comment, and the constants
// below were from http://netlib.sandia.gov/cephes/cmath/sin.c,
// available from http://www.netlib.org/cephes/cmath.tgz.
// The go code is a simplified version of the original C.
//
//      tan.c
//
//      Circular tangent
//
// SYNOPSIS:
//
// double x, y, tan();
// y = tan( x );
//
// DESCRIPTION:
//
// Returns the circular tangent of the radian argument x.
//
// Range reduction is modulo pi/4.  A rational function
//       x + x**3 P(x**2)/Q(x**2)
// is employed in the basic interval [0, pi/4].
//
// ACCURACY:
//                      Relative error:
// arithmetic   domain     # trials      peak         rms
//    DEC      +-1.07e9      44000      4.1e-17     1.0e-17
//    IEEE     +-1.07e9      30000      2.9e-16     8.1e-17
//
// Partial loss of accuracy begins to occur at x = 2**30 = 1.074e9.  The loss
// is not gradual, but jumps suddenly to about 1 part in 10e7.  Results may
// be meaningless for x > 2**49 = 5.6e14.
// [Accuracy loss statement from sin.go comments.]
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov
````

### Original source notice

- go@go1.26.10/src/math/tanh.go (source SHA-256: `a80f128fa85ddc9e94983b53f4a928f6811b96174e7920032e3e1a4367b96591`)

````text
// The original C code, the long comment, and the constants
// below were from http://netlib.sandia.gov/cephes/cmath/sin.c,
// available from http://www.netlib.org/cephes/cmath.tgz.
// The go code is a simplified version of the original C.
//      tanh.c
//
//      Hyperbolic tangent
//
// SYNOPSIS:
//
// double x, y, tanh();
//
// y = tanh( x );
//
// DESCRIPTION:
//
// Returns hyperbolic tangent of argument in the range MINLOG to MAXLOG.
//      MAXLOG = 8.8029691931113054295988e+01 = log(2**127)
//      MINLOG = -8.872283911167299960540e+01 = log(2**-128)
//
// A rational function is used for |x| < 0.625.  The form
// x + x**3 P(x)/Q(x) of Cody & Waite is employed.
// Otherwise,
//      tanh(x) = sinh(x)/cosh(x) = 1  -  2/(exp(2x) + 1).
//
// ACCURACY:
//
//                      Relative error:
// arithmetic   domain     # trials      peak         rms
//    IEEE      -2,2        30000       2.5e-16     5.8e-17
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov
//
````

### Original source notice

- go@go1.26.10/src/math/atan.go (source SHA-256: `777b61ebf72d42b754004fe6a5e9a4ac7fb85fb95cd5e6fef396540ce6c42ee1`)

````text
// The original C code, the long comment, and the constants below were
// from http://netlib.sandia.gov/cephes/cmath/atan.c, available from
// http://www.netlib.org/cephes/cmath.tgz.
// The go code is a version of the original C.
//
// atan.c
// Inverse circular tangent (arctangent)
//
// SYNOPSIS:
// double x, y, atan();
// y = atan( x );
//
// DESCRIPTION:
// Returns radian angle between -pi/2 and +pi/2 whose tangent is x.
//
// Range reduction is from three intervals into the interval from zero to 0.66.
// The approximant uses a rational function of degree 4/5 of the form
// x + x**3 P(x)/Q(x).
//
// ACCURACY:
//                      Relative error:
// arithmetic   domain    # trials  peak     rms
//    DEC       -10, 10   50000     2.4e-17  8.3e-18
//    IEEE      -10, 10   10^6      1.8e-16  5.0e-17
//
// Cephes Math Library Release 2.8:  June, 2000
// Copyright 1984, 1987, 1989, 1992, 2000 by Stephen L. Moshier
//
// The readme file at http://netlib.sandia.gov/cephes/ says:
//    Some software in this archive may be from the book _Methods and
// Programs for Mathematical Functions_ (Prentice-Hall or Simon & Schuster
// International, 1989) or from the Cephes Mathematical Library, a
// commercial product. In either event, it is copyrighted by the author.
// What you see here may be used freely but it comes with no support or
// guarantee.
//
//   The two known misprints in the book are repaired here in the
// source listings for the gamma function and the incomplete beta
// integral.
//
//   Stephen L. Moshier
//   moshier@na-net.ornl.gov
````

### Original source notice

- go@go1.26.10/src/crypto/internal/fips140/aes/aes_generic.go (source SHA-256: `a33d0a61c315c185f168987c7a2f4035d78b1059e9981cf750b7bdf15f91af08`)

````text
// This Go implementation is derived in part from the reference
// ANSI C implementation, which carries the following notice:
//
//	rijndael-alg-fst.c
//
//	@version 3.0 (December 2000)
//
//	Optimised ANSI C code for the Rijndael cipher (now AES)
//
//	@author Vincent Rijmen <vincent.rijmen@esat.kuleuven.ac.be>
//	@author Antoon Bosselaers <antoon.bosselaers@esat.kuleuven.ac.be>
//	@author Paulo Barreto <paulo.barreto@terra.com.br>
//
//	This code is hereby placed in the public domain.
//
//	THIS SOFTWARE IS PROVIDED BY THE AUTHORS ''AS IS'' AND ANY EXPRESS
//	OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
//	WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
//	ARE DISCLAIMED.  IN NO EVENT SHALL THE AUTHORS OR CONTRIBUTORS BE
//	LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
//	CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
//	SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR
//	BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
//	WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE
//	OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
//	EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// See FIPS 197 for specification, and see Daemen and Rijmen's Rijndael submission
// for implementation details.
//	https://csrc.nist.gov/csrc/media/publications/fips/197/final/documents/fips-197.pdf
//	https://csrc.nist.gov/archive/aes/rijndael/Rijndael-ammended.pdf
````


## golang.org/x/crypto@v0.58.0

Source: https://pkg.go.dev/golang.org/x/crypto@v0.58.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:COkYLr4k7nDI4r6QtMML2AfzeIswWHAJthLS6K/M/54=`

### golang.org/x/crypto@v0.58.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### golang.org/x/crypto@v0.58.0/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/md4/md4.go (source SHA-256: `9be73f76d4488c7a216828e4caec1d2701955f2b03b77abd019725b0442e651e`)
- golang.org/x/crypto@v0.58.0/md4/md4block.go (source SHA-256: `0303098d7fe87ad0ce4b8302dc984e14f1de47c53f071bbf97b808e7c81575f8`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/blowfish/block.go (source SHA-256: `4f7eecb319e7c3d4018c0f8528be82cd889ae7b2829ee1f9c1f83a651a883b07`)
- golang.org/x/crypto@v0.58.0/blowfish/cipher.go (source SHA-256: `c67e62fdf539aa94cd50c1055bdc1d4c6dc93f4b2d654addfd97ccd1d12f6e48`)
- golang.org/x/crypto@v0.58.0/blowfish/const.go (source SHA-256: `fbf914c515f87da351d3b8d1ffb3751ad73a26f6817ed0857470bd44c9479c0d`)

````text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/ssh/channel.go (source SHA-256: `b753096dc6bdec0533c808bef0905663b72745d2739ed78bd3e0e6653518326a`)
- golang.org/x/crypto@v0.58.0/ssh/cipher.go (source SHA-256: `e7aa62253f84266bf5a90a8d26f908fa36cbf9732a09d775bc83e0922bd1fca6`)
- golang.org/x/crypto@v0.58.0/ssh/client.go (source SHA-256: `614d95be492c97e247071ce22bc2d0f53ae206944814b9791227c6435cb3fd67`)
- golang.org/x/crypto@v0.58.0/ssh/client_auth.go (source SHA-256: `f26b065303604d257be5382138a4c2f6b707911653353811edc3811b8ce3b153`)
- golang.org/x/crypto@v0.58.0/ssh/common.go (source SHA-256: `506eead19076f97e40400e9e52a4b37caf8f8c468d5e0fa97f341eb0616b9ca8`)
- golang.org/x/crypto@v0.58.0/ssh/doc.go (source SHA-256: `5c806706f6169e16f69ffc44fe1723fe1674c040280364a8b931b8abf2f77a32`)
- golang.org/x/crypto@v0.58.0/ssh/messages.go (source SHA-256: `fbeeb5fffa18f017609787d1ccf9d177c1d44f65ad784add5680e51000161c45`)
- golang.org/x/crypto@v0.58.0/ssh/server.go (source SHA-256: `a9ff150f87775be31f40685d38f1c6934101a461f9bce82936a35331f6949c82`)
- golang.org/x/crypto@v0.58.0/ssh/session.go (source SHA-256: `6358e49e2753b7799751e9ad7cd3b71eb2dcb5916602e39aac70ef40d0947b0b`)
- golang.org/x/crypto@v0.58.0/ssh/ssh_gss.go (source SHA-256: `910430ffe74240b1b9d9689a6e7b4c72e1a90faff284aceae5fdaeff24ade183`)
- golang.org/x/crypto@v0.58.0/ssh/tcpip.go (source SHA-256: `d59df8f49249d436b12408f54b4bd51a229ae52bdf067ed26669697075700e6f`)
- golang.org/x/crypto@v0.58.0/ssh/transport.go (source SHA-256: `245b42435040ca899da4ca5e3b1022ab93e1dfdfbdf3cbbac06d4c656b7f344c`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/internal/poly1305/poly1305.go (source SHA-256: `c4b20ab7330e47ee155a775f69ae3b0e432c75d74ad005cf764228a69ddf6da6`)
- golang.org/x/crypto@v0.58.0/internal/poly1305/sum_asm.go (source SHA-256: `39e3031ebe8eccf3ac0324c1d8d21fbbac8b74e9bd48737b0c042c28b4676d43`)
- golang.org/x/crypto@v0.58.0/ssh/buffer.go (source SHA-256: `8a7997c253d676a137735b1d0d9b17b8bca33dc54d351bd86620a1e017ad3520`)
- golang.org/x/crypto@v0.58.0/ssh/certs.go (source SHA-256: `093b1448a7e8a17a5056e31e007bd10cda2e0a67a74b8548cede45068b4b1f2c`)
- golang.org/x/crypto@v0.58.0/ssh/keys.go (source SHA-256: `805c102eb987b391c6384fa3d4c31a025c7821e70a1a703d3d6443745881edb3`)
- golang.org/x/crypto@v0.58.0/ssh/mac.go (source SHA-256: `0894452248f590b8b26b97fa38b86116f9ef86305b4404096e4103e40537161c`)

````text
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/ssh/connection.go (source SHA-256: `57d6a5210afca54f1ff1c39b226bb3b11e035484c0e607bbf3934eed69747604`)
- golang.org/x/crypto@v0.58.0/ssh/handshake.go (source SHA-256: `2eea5931af6dab0bde843654eea304b9ba0cf0a4a579870be822c1725ec6d95d`)
- golang.org/x/crypto@v0.58.0/ssh/kex.go (source SHA-256: `0aed5134fc6e67e02fd40eca62a70547dee4d9317dffb3eff0e35052b5c39661`)
- golang.org/x/crypto@v0.58.0/ssh/mux.go (source SHA-256: `9c560d448d8836de7a929efd604a7e7f863750d50675bfdd6473400f4b4fb96c`)

````text
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/hkdf/hkdf.go (source SHA-256: `4f0e42b70b0cebbe27e9092072c48e579047fe74767b3d9ed7c35349d57c0e5e`)
- golang.org/x/crypto@v0.58.0/ssh/internal/bcrypt_pbkdf/bcrypt_pbkdf.go (source SHA-256: `4694b500beb45c7f61c050b59bd7f7aef88f87790234fa49d78379b032d63d9a`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/chacha20/chacha_generic.go (source SHA-256: `34403e82b1387b4402b00ce30c1364508c333f4bdbe671321690c6ebaa8d3180`)
- golang.org/x/crypto@v0.58.0/chacha20poly1305/chacha20poly1305.go (source SHA-256: `08400b984074331e97c7e2a1d9701f14cb8f9919b3eb3e0d1bb8e3549961cac2`)
- golang.org/x/crypto@v0.58.0/chacha20poly1305/chacha20poly1305_amd64.go (source SHA-256: `bae40d20a6077f0f5570b28086cdc4cb6e51aaa757e54b32ba99770f8b234869`)
- golang.org/x/crypto@v0.58.0/chacha20poly1305/chacha20poly1305_generic.go (source SHA-256: `f3e04c8517e163f94ec2bb6e7b634ef239432ed7422b1d8695cfc446c8551181`)
- golang.org/x/crypto@v0.58.0/chacha20poly1305/chacha20poly1305_noasm.go (source SHA-256: `0ff7a84348991f03aa8c9a74472b112036cb8149218aa284ed034c2ea011afc6`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/cryptobyte/asn1.go (source SHA-256: `a7b568d1a846cfccee16f7847c5a0dcc5326c3c36dd299cc69956d9c1ef164bc`)
- golang.org/x/crypto@v0.58.0/cryptobyte/asn1/asn1.go (source SHA-256: `306307e55077ec60b0c90dd90e1a7954805349ce2fd68d66f08a660524653f80`)
- golang.org/x/crypto@v0.58.0/cryptobyte/builder.go (source SHA-256: `63e388880001b25f0b8bceb2b3702699146f3691f35dd5284a7651aff93fcf1c`)
- golang.org/x/crypto@v0.58.0/cryptobyte/string.go (source SHA-256: `66485c8123b338056d97077938776955fe6c080e05c2c77e36501946113d68e9`)
- golang.org/x/crypto@v0.58.0/ssh/knownhosts/knownhosts.go (source SHA-256: `be5cf474f045ed8d912daa50b69ddd46122f3ddfc5ee2a67bd7b560752f28357`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/chacha20/chacha_arm64.go (source SHA-256: `50974d7653c355d8356aaa33318f1df5d9127b5dc96a38b5d88e91aeb188feb4`)
- golang.org/x/crypto@v0.58.0/chacha20/chacha_noasm.go (source SHA-256: `7533b8266ccb098f4cff5649dddff8a58574411c4765712d5bb045484bd078b1`)
- golang.org/x/crypto@v0.58.0/chacha20poly1305/xchacha20poly1305.go (source SHA-256: `7325ab1f70b6ab4a241ca0237f050202ddb35225a690c8146006e51a6ceeb685`)
- golang.org/x/crypto@v0.58.0/internal/alias/alias.go (source SHA-256: `7bb0ef254a7e9360ec43e3129fad11d27edcc611bf87c12dba089f6e381f462e`)
- golang.org/x/crypto@v0.58.0/internal/poly1305/mac_noasm.go (source SHA-256: `e8ad079a9c995609bb5bf27fe41d597c5a0bda8ee0b887cb6f23754427c107ea`)
- golang.org/x/crypto@v0.58.0/internal/poly1305/sum_generic.go (source SHA-256: `b0094a2895d5bda42dcaaf57c0b31fc914c3b6c6aa6237aab0100a6d78346933`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/chacha20/chacha_arm64.s (source SHA-256: `73036d54d1961a9e5cb8ce57298bd7803443dbe0de7cae4a9143b62d3a7d84b5`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/chacha20/xor.go (source SHA-256: `c3adc8555766e36629b4dd3fcd579969d2d2737c65a379db997253a2b7b18072`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found src the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/curve25519/curve25519.go (source SHA-256: `af22c7f5967f190e56031c2d98e750f9dd55788cbfa4c7a2b9cd5321f9cf3735`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/ssh/mlkem.go (source SHA-256: `0e78c5f3f901d309aab76b40d83d3dc473908d42b4ddf16d67ea72a5e3fd2c35`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/chacha20poly1305/fips140only_go1.26.go (source SHA-256: `60368cae6d4d630509309819d7f397532c6f7812978c12983bca375e4ae810aa`)

````text
// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/crypto@v0.58.0/ssh/control.go (source SHA-256: `a0081d8573ea786d692228b09e530c485234a98bd300188f6f94495530293357`)
- golang.org/x/crypto@v0.58.0/ssh/mldsa_unsupported.go (source SHA-256: `fcceb8c1c620cea24f450670e590a24bce2dfdaae84485e1f94ed34890368daf`)

````text
// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## golang.org/x/image@v0.47.0

Source: https://pkg.go.dev/golang.org/x/image@v0.47.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:za/8gVz++PPZqmZJBNSISICN5Rrj5zixP/8K++fgmqo=`

### golang.org/x/image@v0.47.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### golang.org/x/image@v0.47.0/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- golang.org/x/image@v0.47.0/tiff/buffer.go (source SHA-256: `0acd2ec75604a09d527f17c3efffef34ae5f3e1403b28455d24f6f8c23456f01`)
- golang.org/x/image@v0.47.0/tiff/compress.go (source SHA-256: `33b980399c8a266c18560db9a2f01836bdb2a74c46f02c8da63acd1d569593e9`)
- golang.org/x/image@v0.47.0/tiff/consts.go (source SHA-256: `700f8b3a05bcbc5a7d12f1dafa26f049ab9f29926fe4b951fa15dfd32f4d5a53`)
- golang.org/x/image@v0.47.0/tiff/lzw/reader.go (source SHA-256: `978237289711b56c7d5b26c773a84ca62126687ce846fc699906ac4cc30bdf8f`)
- golang.org/x/image@v0.47.0/tiff/reader.go (source SHA-256: `11d080e5d7f2d9961a49f89412240e3f310d0edcee346ccfbd463a2596915803`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/image@v0.47.0/tiff/writer.go (source SHA-256: `62de278c4b115c002dc7890b5efc49e9bd0e192252b48007553cb9b273ac78ee`)

````text
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/image@v0.47.0/colornames/colornames.go (source SHA-256: `d30e7c6a5b3eb491b6e119eae75e00ca7a5d28fcd2679de5105e8868f26f8346`)
- golang.org/x/image@v0.47.0/draw/draw.go (source SHA-256: `78ea13baad5e523f1bc2e213848e46622a22a1324c9e8cf4a944ee0c85ab168d`)
- golang.org/x/image@v0.47.0/draw/scale.go (source SHA-256: `c38c400ca8f00628f0eae1f65c8cba81084ae2836c8575b609e38484a06d6f92`)
- golang.org/x/image@v0.47.0/math/f64/f64.go (source SHA-256: `943538d75254f50e03d00812a4546d07e9e305caad1b56ff52d61e66b90fd439`)
- golang.org/x/image@v0.47.0/math/fixed/fixed.go (source SHA-256: `b74f2166fce03cecfebb00d3e9641b559dab8ffb00d72139c1dc69ed1ac05725`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/image@v0.47.0/vector/acc_amd64.go (source SHA-256: `eb6b485c56b5cd88cf6946b2ea67d21fe1deac48b89113a94e39965f147bd49a`)
- golang.org/x/image@v0.47.0/vector/acc_other.go (source SHA-256: `ef5e973825a5ead655e348011061ef2f3a2ce50cb21f535af9bbd3816c69e479`)
- golang.org/x/image@v0.47.0/vector/raster_fixed.go (source SHA-256: `c80fc497b5650b7b62f56452a31ff818f7fa78497280b7c999ed80d479d63cc8`)
- golang.org/x/image@v0.47.0/vector/raster_floating.go (source SHA-256: `7dce1afc53db87bb0ada484db9871b84e352c39e34e6caf7b29d0a65ef9770a9`)
- golang.org/x/image@v0.47.0/vector/vector.go (source SHA-256: `f908264b582ab128d16632b8a73600f63cf3a60e8078cf4b9662dc8087e438d0`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/image@v0.47.0/ccitt/reader.go (source SHA-256: `0e4fbd770117cb76edb0e43573a5f6b74d979a850c873ae10df44c2c03cc5178`)
- golang.org/x/image@v0.47.0/ccitt/writer.go (source SHA-256: `9a5d232614aa458070fc0b05ede57cb87cac7d715283509caa4bdfbe83762cb0`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/image@v0.47.0/internal/safemath/safemath.go (source SHA-256: `0cfc3a6086a878ada0ce02d92e8995011346e9e8b1a2ca9d959b2c2c467cd327`)

````text
// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## golang.org/x/net@v0.61.0

Source: https://pkg.go.dev/golang.org/x/net@v0.61.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:RR+6j/BTBCrggOLugxa+66R8FvoN6MHlzDlZIbN9YFw=`

### golang.org/x/net@v0.61.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### golang.org/x/net@v0.61.0/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- golang.org/x/net@v0.61.0/html/doc.go (source SHA-256: `e5d486a3b703d17efdf30cc9b3af6649b262c89e74cd4de285e0f079797b3c6f`)
- golang.org/x/net@v0.61.0/html/entity.go (source SHA-256: `680a9bb55e7eca9bc0cbee59755591ab9f97abf3fe39490211eb088cc2cd62c2`)
- golang.org/x/net@v0.61.0/html/escape.go (source SHA-256: `d0dfaa8c3b36ef9d4578dfdf9480c654252afc6a09128d96d570e74c563c82c9`)
- golang.org/x/net@v0.61.0/html/parse.go (source SHA-256: `ab68a3276498d29551d38faebd93a622c4d1d908f6f1dad00bce438e74110726`)
- golang.org/x/net@v0.61.0/html/token.go (source SHA-256: `28e4cab66406c656cd2ae007517bac998b3d68d18fce0e0717ed98bed44ed336`)

````text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/html/const.go (source SHA-256: `44c814afac4b0206a4f2af92607f93a72b93919ec2c3884d372edcdc24507584`)
- golang.org/x/net@v0.61.0/html/doctype.go (source SHA-256: `b629b17ed04bcab8d832d8caaee00f458c180fd31b6838e318a439a4bf87fffe`)
- golang.org/x/net@v0.61.0/html/foreign.go (source SHA-256: `d227ce6839726bb85fc88ec95df265901da8a04d7e2f7261101bdc7b4d1af1ad`)
- golang.org/x/net@v0.61.0/html/node.go (source SHA-256: `3acf739c450f61abd0fc45ccb19e4ddde2982b798791037a6a03c78603b6176f`)
- golang.org/x/net@v0.61.0/html/render.go (source SHA-256: `c53cd811db01490a88a3932ad97e0d0622fa87a8533662fb9810bc06814c4af8`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/html/atom/atom.go (source SHA-256: `86490b2027af4443dacfe96c313a6c632c03ea1acdf85f7f61c4fdde32bc2f3d`)

````text
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/html/charset/charset.go (source SHA-256: `7810cb77bad4b115611626e3d0415126419133c5d6bb9c7424b45d721a7afdc3`)

````text
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/databuffer.go (source SHA-256: `af910b5796864e45f849321fd18b71bd4ab420ee239f5e9274288e3e0b8541de`)
- golang.org/x/net@v0.61.0/http2/errors.go (source SHA-256: `fde1284c85d1730c86183ce604cfb22e917edfcd4ed9544a4a12d511befbea49`)
- golang.org/x/net@v0.61.0/http2/flow.go (source SHA-256: `fedf89792079bcd134e4c41da64f5ac99a0881a681cf5b642c9e086d4e655250`)
- golang.org/x/net@v0.61.0/http2/frame.go (source SHA-256: `4466f8ba7f851901df67721935d75f526f9f91184bd93c899307e577b97d9bd7`)
- golang.org/x/net@v0.61.0/http2/gotrack.go (source SHA-256: `2e404e35aa3d9c0350ce72191e5b1b525c6bb8ead1ce5bf8d8533aaf1c8efd3c`)
- golang.org/x/net@v0.61.0/http2/hpack/encode.go (source SHA-256: `b39a13dcfd94e03231ab39fe0cfdde3695a6674a481575385237f34dc669d846`)
- golang.org/x/net@v0.61.0/http2/hpack/hpack.go (source SHA-256: `e1a1b7c22e9df2d1d2da30ac382347a9076ba8736084635957bc7a6d7ddb5034`)
- golang.org/x/net@v0.61.0/http2/hpack/huffman.go (source SHA-256: `3d41693bd77da34093c1d334842d8cbeffb3a2f3e0fa5cf0af5e6e9e65b19be1`)
- golang.org/x/net@v0.61.0/http2/hpack/tables.go (source SHA-256: `adfd7d2e7b49b99ab21ae3ce296c8d9d6bc657362e6b5e522fe89065d67bd48b`)
- golang.org/x/net@v0.61.0/http2/http2.go (source SHA-256: `e1b38f7628feed3fd9482d2e705f07a8aecf58d8068c09e012a56b6b448452c2`)
- golang.org/x/net@v0.61.0/http2/pipe.go (source SHA-256: `3fbe00025d441802c154549f5bd7562b9937ff43f1ce64c987bda75fcaf7f518`)
- golang.org/x/net@v0.61.0/http2/server.go (source SHA-256: `28910755e3963370c2eb304c7f3d2c501df004e8c67b69969a0675047011260d`)
- golang.org/x/net@v0.61.0/http2/write.go (source SHA-256: `be8997734245f7b15c72cbc65f6ed20ebef6d7ae463a14d733e777c20661ce5f`)
- golang.org/x/net@v0.61.0/http2/writesched.go (source SHA-256: `1c7a6a14d8658a6c41547c79dc36c56fafe73a4544fb77825a7663c5b40fd7fd`)
- golang.org/x/net@v0.61.0/http2/writesched_random.go (source SHA-256: `071126c1cc5308dd8b37c01dcf9bab2b8cdcb3ab4fc237b03740b8d257ecf873`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/client_conn_pool.go (source SHA-256: `ee773ca6b4e00a999f916214138ca74ea06a9d3d8d6fdffba354871587019dcf`)
- golang.org/x/net@v0.61.0/http2/transport.go (source SHA-256: `65813b29f50cca84bfe9b1df9394662f40417f80653c53f8a97397e9aaff6e35`)
- golang.org/x/net@v0.61.0/internal/timeseries/timeseries.go (source SHA-256: `fb2a98dacbeb1a456b79bf48d4c90ec1f213fad4cbafd87564326b61e5d909ce`)
- golang.org/x/net@v0.61.0/trace/events.go (source SHA-256: `b0ab4685a199b6246ff459b7c7ac0e1d47dfa6983d72efc11b2948c9f99cd674`)
- golang.org/x/net@v0.61.0/trace/histogram.go (source SHA-256: `7110d26994a380e7e2d06095672a6d1e1fe006852f0e144938f7cc24486c961e`)
- golang.org/x/net@v0.61.0/trace/trace.go (source SHA-256: `be54020d8d18797b9a9d727b68483495106bbe211b528db42dae6636b71d1fc7`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http/httpguts/httplex.go (source SHA-256: `8cfd5819dc821d45f82fde10afff4b2e652535b0091e0559fc8ad7766051a2c2`)
- golang.org/x/net@v0.61.0/http2/writesched_priority_rfc7540.go (source SHA-256: `423a3617111d0b5d44d1e1c5718f89e590b07eca56f48b41073b79a7f956ed77`)
- golang.org/x/net@v0.61.0/idna/idna.go (source SHA-256: `26de47c66ae70e53c5bedffca7fd61aefcde2ef8df6db37c9a50a3fb3c152aab`)
- golang.org/x/net@v0.61.0/idna/punycode.go (source SHA-256: `3e65858245f1a7d32e225aba714a305d3ef34496acb839956a7b15d1c7a34eb7`)
- golang.org/x/net@v0.61.0/idna/trie.go (source SHA-256: `8d2a16a56e4cab9d23a7afb49443bfe7a5ff0d74e885acc9121b4a9c598825eb`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/ciphers.go (source SHA-256: `623262af28491cd81b3570936cd06757f629d2ae4fb6ef087af6762f30d7a5eb`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http/httpguts/guts.go (source SHA-256: `3da54270a625540809f2f767e54e18acf27199aca64b6cf16d6d3fa3999de959`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/ascii.go (source SHA-256: `38ae63dd18385b2495931691e043a77ba73f303bc2f742d3a52dafacd306e8b4`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/writesched_roundrobin.go (source SHA-256: `1168ebd908b1cae2e695e76babde24cd33a8d085ceaf93fe11c1582de7570a2b`)

````text
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/html/iter.go (source SHA-256: `54299c163639df83f3dd22ff31643b6b6cedbedb0c388b1f53020146a10d7504`)
- golang.org/x/net@v0.61.0/http2/config.go (source SHA-256: `d167984b5523fecaf2b99b626c7a4be2d7fc2268b6719b01913263cb80ca4d5a`)
- golang.org/x/net@v0.61.0/http2/unencrypted.go (source SHA-256: `8efbd1156e499e01d14dd2c8095a5318256d34f24d25fd74ced760dde5b5ab42`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/config_go126.go (source SHA-256: `59a922ba01b0c7a4c2a5155613306f2bce7fb358992541b321b6df5c594dd2c6`)
- golang.org/x/net@v0.61.0/http2/writesched_priority_rfc9218.go (source SHA-256: `cabdfd7cd04d6723d11b5f26b39ed97517e99a665896cf26f65cbabdbb761e86`)
- golang.org/x/net@v0.61.0/internal/httpcommon/ascii.go (source SHA-256: `7446cfcfe70f682d5ad3102ca0aeef3bfbef86c799f0c61af9be9e942fe950dd`)
- golang.org/x/net@v0.61.0/internal/httpcommon/headermap.go (source SHA-256: `ee95dbc37a39a5a9ac3bc048ba3233f13d3f5e18e2faa8ba44d588218a888ade`)
- golang.org/x/net@v0.61.0/internal/httpcommon/request.go (source SHA-256: `90095651af5ec9560583099b5783bed12c17f23ffee19c3adfee83785d992fa0`)
- golang.org/x/net@v0.61.0/internal/httpsfv/httpsfv.go (source SHA-256: `bb64129bf25faa9065a029ec0f0b7ffd61abbb49e561100ff66614ab5705001a`)

````text
// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/net@v0.61.0/http2/client_priority_go126.go (source SHA-256: `b51147502bb7e74c384f1dd249c951920055b380d18fe38862532f8f24736db9`)
- golang.org/x/net@v0.61.0/http2/clientconn.go (source SHA-256: `75517095356686258410eebac3a058152f4e4fcf85d98732e55e72db3d8860b8`)
- golang.org/x/net@v0.61.0/http2/server_common.go (source SHA-256: `2fb45d3049aa939625c2bfc8e4faa546885a6f25311105f08151fb35710769e8`)
- golang.org/x/net@v0.61.0/http2/transport_common.go (source SHA-256: `63fd3747e63a2b2ad9de77ed12b2665d672ebde1a35cc840179cf98e50a4cdb6`)
- golang.org/x/net@v0.61.0/http2/writesched_common.go (source SHA-256: `fed2317a546f28aef250b55971a44621218c4076a30e03bb57369839998ba202`)
- golang.org/x/net@v0.61.0/internal/httpcommon/gzip.go (source SHA-256: `8ed26b633263782a1f18b9ff71e420fa54777bc61078989157c9dd90f69455e8`)

````text
// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## golang.org/x/oauth2@v0.37.0

Source: https://pkg.go.dev/golang.org/x/oauth2@v0.37.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:JUlcxA8oAtauLfiH8FX2/FkAWHAdi0QtGCGc+hofE98=`

### golang.org/x/oauth2@v0.37.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/appengine.go (source SHA-256: `def4457533b463cffa49c657cbdf721814aaf90b1c44217e2b96bbb8a10350ce`)
- golang.org/x/oauth2@v0.37.0/google/google.go (source SHA-256: `81f3e0573ea4fc3f8132c5ed1cff52b261343e1ffaa0dff773b5e8eb94ce01ed`)
- golang.org/x/oauth2@v0.37.0/internal/oauth2.go (source SHA-256: `d2ce47d21fac46b05aed8922df22b0ffb8ee3896b1198a0ba6de669224288e1a`)
- golang.org/x/oauth2@v0.37.0/internal/token.go (source SHA-256: `98972e4d8b7a8807d2dc6ad760ad4bac346391b2d0bb155d6a7c62bfe6844a22`)
- golang.org/x/oauth2@v0.37.0/internal/transport.go (source SHA-256: `20900bde7fa0057d3f688547754e664bd41497d76355c8eb79c30896cc842d47`)
- golang.org/x/oauth2@v0.37.0/jws/jws.go (source SHA-256: `9034dcb95fae39e25dfa3dcc1e4541a6a30241cf340744d4e63d1d97c6d94401`)
- golang.org/x/oauth2@v0.37.0/jwt/jwt.go (source SHA-256: `279de50daad4b07dcabba0482ef69c6873230819454e4a71605481464f61088c`)
- golang.org/x/oauth2@v0.37.0/oauth2.go (source SHA-256: `e6ccbbe51f721b73c65fe60597b024e01c2908353fb4d03919fd8899d3d3f968`)
- golang.org/x/oauth2@v0.37.0/token.go (source SHA-256: `9715ac669ec1d80ae4549dd83379a95c7ea278542a9c96ad55e957e2e822c8db`)
- golang.org/x/oauth2@v0.37.0/transport.go (source SHA-256: `27f9235d8558e77570a7845a5614f871cbd96b235c6d208332fa1f081e846191`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/default.go (source SHA-256: `a0e67a6d6542e90e45af653b4c41681973dec5ce54e447cfd524f59227a31c6d`)
- golang.org/x/oauth2@v0.37.0/google/jwt.go (source SHA-256: `0097a6e01b047951affe253f5b4760aa79bf496ccd468d4c2f5b970660efb6c1`)
- golang.org/x/oauth2@v0.37.0/google/sdk.go (source SHA-256: `5a6023c614294ed08ce271396a7a6f46fe682ed02afb6ba646555cb250cc8cf1`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/internal/doc.go (source SHA-256: `4c268b85467124f37ca0ff613a54fd07ba33425dfb53423fda209779ecae3f2f`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/doc.go (source SHA-256: `d3382d4a7896d46ce2869301123db04ae90f373ccf41e153a042d931c684b05b`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/externalaccount/basecredentials.go (source SHA-256: `d092176c6fde8344de456548d5b0db785b1abff546f13b2a880fd494f501d579`)
- golang.org/x/oauth2@v0.37.0/google/externalaccount/filecredsource.go (source SHA-256: `c32927a2e5b4627ae196bbfc393059192980d109d96a3b4d423ff86a9fc354ff`)
- golang.org/x/oauth2@v0.37.0/google/externalaccount/urlcredsource.go (source SHA-256: `2fe79ec7e62770032591e608a62b088e6dbd1d30d46170cb0786d87e36643dbc`)
- golang.org/x/oauth2@v0.37.0/google/internal/stsexchange/clientauth.go (source SHA-256: `ec778b10353449d358aa3d3a8525c90aae35f148d7f62ac9c0d0fcd3c4ec92c1`)
- golang.org/x/oauth2@v0.37.0/google/internal/stsexchange/sts_exchange.go (source SHA-256: `1a535e3acb5813680e9fb8f277bb0db212898704cf7781d11ec2e21d62cd3f47`)

````text
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/authhandler/authhandler.go (source SHA-256: `983e143df172549fc95158e06d0e8ade55300fcd67066793edab18b5757d4002`)
- golang.org/x/oauth2@v0.37.0/google/externalaccount/aws.go (source SHA-256: `9af6d7da0fde7a1a19858de8f5629891961ef878192c76e47a18792f0cdf3c81`)
- golang.org/x/oauth2@v0.37.0/google/internal/impersonate/impersonate.go (source SHA-256: `65c7a38c4ad49f859269ad21a5b7f6c4ea6c85a69ae77dc9c5df60b9d8e521d2`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/error.go (source SHA-256: `e2f0a6167343d0e8703aa6aafb7ee9fe662d0834e9d8fec455811628a43d26ab`)
- golang.org/x/oauth2@v0.37.0/google/externalaccount/executablecredsource.go (source SHA-256: `055223607710f94c6d4ed25fd601f4d3e7d40d2b6c6d850d6f97b8bb5bb8d20d`)

````text
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/externalaccount/header.go (source SHA-256: `3086bd482007eec1cfe6b6898f84ab3ee777bf7d2f14861148351a1a9cc0d3db`)
- golang.org/x/oauth2@v0.37.0/google/internal/externalaccountauthorizeduser/externalaccountauthorizeduser.go (source SHA-256: `6b18306bcff0398afb3d2d8f74d7c4f2f11a91a8539f267bb86f6c97408e58e0`)
- golang.org/x/oauth2@v0.37.0/pkce.go (source SHA-256: `865223cbb516af7142480aff00be0c7736938baa460c4f541598a2e9464f885a`)

````text
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/oauth2@v0.37.0/google/externalaccount/programmaticrefreshcredsource.go (source SHA-256: `e24fb19204374c57cff9dc3beb1506176529e190d1b56eabdfe9541d6fdb0d78`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## golang.org/x/sys@v0.49.0

Source: https://pkg.go.dev/golang.org/x/sys@v0.49.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:XbzkgYJHdqh/8m2Uu0W/dQv8nktxx4BFHp1M0gROTXA=`

### golang.org/x/sys@v0.49.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### golang.org/x/sys@v0.49.0/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/dirent.go (source SHA-256: `03e3b15a8428e2f1520386291052fd30e5d74ecb4d78c724bae7953d71425be1`)
- golang.org/x/sys@v0.49.0/unix/syscall.go (source SHA-256: `2f1f25726d54e971b2dcedc7e8bc7d109d579e3a671f0ecfabd4deb85246c3be`)
- golang.org/x/sys@v0.49.0/unix/syscall_bsd.go (source SHA-256: `9ae0c12d64f207f6f9cee249a20cce822595688d99440363b673ec798cf10847`)
- golang.org/x/sys@v0.49.0/unix/syscall_darwin_amd64.go (source SHA-256: `7a180893c051cc50d668414ce2ae60da968b7303aa3b1ebb816fcb0619b7faa5`)
- golang.org/x/sys@v0.49.0/unix/syscall_linux.go (source SHA-256: `5cba74fe971377ace2b1192f8d7818f991f3fd8b11e15639960c93b1c5771220`)
- golang.org/x/sys@v0.49.0/unix/syscall_linux_amd64.go (source SHA-256: `dfb1ba340e05d035b27aec5624c4743b5b686607ece9ab8d39eed329ac179df4`)
- golang.org/x/sys@v0.49.0/unix/syscall_unix.go (source SHA-256: `d851dcf05549674486f35d58b0357bb5c1ea9378bd234b1646534a35dc4a6da5`)
- golang.org/x/sys@v0.49.0/windows/exec_windows.go (source SHA-256: `9c4824d3c8a1d54189f9d81e54a5786971b338c6d8b116ad188c1746212534a8`)
- golang.org/x/sys@v0.49.0/windows/str.go (source SHA-256: `6c8cc24a3e0c11c6b4950dbb49a5f7eb72e65fcb861a8cd71cdad376125a43cb`)
- golang.org/x/sys@v0.49.0/windows/syscall.go (source SHA-256: `60b6c688bece9fb3550cb369311b55dc33d1e2350a93f9c8a5f8fa0cd8ddb2cc`)
- golang.org/x/sys@v0.49.0/windows/syscall_windows.go (source SHA-256: `589a6e28bca016f69eb25b75e82a808162d707de98d3d50367337303dcdd7018`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/asm_linux_amd64.s (source SHA-256: `14c826e5d2db337e49c32e0b5a66317b58da198874a0eb950c33aac571e9573c`)

````text
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/syscall_darwin.go (source SHA-256: `2a4add0c5f627a675f7680b6bb36d7f090f2ed1c67bc59f2a0ab882deedfcd75`)

````text
// Copyright 2009,2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/env_unix.go (source SHA-256: `bcb73ccfc5a8dae1f59e4debba69e0f600155b947f358c939bc753443a7a8007`)
- golang.org/x/sys@v0.49.0/windows/env_windows.go (source SHA-256: `fcfa69bb220a15cc151240ecc57bb144cf25c6628e899bc72d8217171376ab5b`)

````text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/sockcmsg_linux.go (source SHA-256: `51d948aa226aa043c2760288aa883f9dbd53b9d65cf8bc20215146fc016dde6f`)
- golang.org/x/sys@v0.49.0/unix/sockcmsg_unix.go (source SHA-256: `7863093bfa7ec9e564c36992fbf2b7c97244c84d406e5fdc60b72c56a23af396`)
- golang.org/x/sys@v0.49.0/windows/dll_windows.go (source SHA-256: `1c62d9345dd0c1d8e660792c1955a407b2049210c3bd2d8000004375f07afed9`)
- golang.org/x/sys@v0.49.0/windows/types_windows.go (source SHA-256: `97dfc2cc0606e42b3bec07fb8c669098cff7ac1a25441df4aafef27dc40fc886`)
- golang.org/x/sys@v0.49.0/windows/types_windows_amd64.go (source SHA-256: `c375406875cabc6447a9a8ca917eccab2227ea33a17ac3ed02a3a8fca0af8ce1`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/race0.go (source SHA-256: `8a78192af2b20cd177f78e3761e10d66567373cef3b7671aa3c7912d65de4c90`)
- golang.org/x/sys@v0.49.0/windows/eventlog.go (source SHA-256: `0162de9581d94b7f73c2125e310a8c08d1f5f6ec9d99d68b450ea330cdab33e0`)
- golang.org/x/sys@v0.49.0/windows/race0.go (source SHA-256: `fd656479b0761b07939c12b95733002fe0437f33e6a9b317cdaf8b3fe3e116e4`)
- golang.org/x/sys@v0.49.0/windows/security_windows.go (source SHA-256: `a81ceadb4023d6712587d1d9d0b5a51948dfbfb8a5786811d3020fe55182b6c6`)
- golang.org/x/sys@v0.49.0/windows/service.go (source SHA-256: `88b18aaac8358e790f5f723acac90799e84b813bc02c9c7ddad2e0f15a0503d7`)

````text
// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/fcntl.go (source SHA-256: `e8d84df0c6ed38a56014cc73d431aa973f1d298169735eaf9f7d42cbb26f7768`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/constants.go (source SHA-256: `f3405abc7484992964143eac589b951132d8e3a90f8359fc1a6e9ddcc201aa8a`)
- golang.org/x/sys@v0.49.0/unix/syscall_darwin_arm64.go (source SHA-256: `ca9f22a90e5e81a736c0ad500f015b139cbb05a23ffdbc39906c6feb6975244e`)
- golang.org/x/sys@v0.49.0/windows/registry/key.go (source SHA-256: `59aab3af1dd13c55cb65c9189cc3da277fecfe82dc4c21baf557cb1c3d2c128d`)
- golang.org/x/sys@v0.49.0/windows/registry/syscall.go (source SHA-256: `c630e393502f67d3c8369774ba3b3a6b09df646e22b788e7298237d3b38c21ab`)
- golang.org/x/sys@v0.49.0/windows/registry/value.go (source SHA-256: `abf23b1dc18dbb3d54ea1e113794cf6d564c2dd1f946e7159f9aadabb8a1a1a5`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/bluetooth_linux.go (source SHA-256: `72ce84a68d3647cd34f04e9040861333ac9e9f9ac5b00008d41514c79da664d8`)
- golang.org/x/sys@v0.49.0/unix/syscall_linux_amd64_gc.go (source SHA-256: `6c5b86f3c4aa3dcf8619c5d444491414f72e014767e57fd9ad62f57653f4d837`)
- golang.org/x/sys@v0.49.0/unix/syscall_unix_gc.go (source SHA-256: `2d2712c417a71c258783dde275c8f9e07c16754c19cd035fa4e3c76def80169a`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/endian_little.go (source SHA-256: `bc06276262c57cf21e35c13ff2f5fdc96c18feb32a055c6201edfa76e95c967f`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
//go:build 386 || amd64 || amd64p32 || alpha || arm || arm64 || loong64 || mipsle || mips64le || mips64p32le || nios2 || ppc64le || riscv || riscv64 || sh
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/dev_darwin.go (source SHA-256: `9a0bc8af77b4325bb10b651e00b8f7974cc972d0e5456a370f2c46a56181ada7`)
- golang.org/x/sys@v0.49.0/unix/dev_linux.go (source SHA-256: `eec381cf025f58965728544dae5df9394d6bebc7e68572e429814c87ee6dd284`)
- golang.org/x/sys@v0.49.0/unix/pagesize_unix.go (source SHA-256: `5fea300c32898efe55341abb39489aeac76f81ea9df513f68d1b2556be7b48a3`)
- golang.org/x/sys@v0.49.0/unix/timestruct.go (source SHA-256: `d0d07c2481ce2692f4e5728d65ccab1d105ecf64f52ac6b9a923ef106088a249`)
- golang.org/x/sys@v0.49.0/windows/memory_windows.go (source SHA-256: `070342979027b739e224a8d35c5dabd9ac6ea662b9a58954446f27103c2443c4`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/cpu.go (source SHA-256: `b3ec218edc560fbb89ff9b5049dcada60b0aac8ff7cf31c98155081ff4a26d88`)
- golang.org/x/sys@v0.49.0/cpu/cpu_gc_x86.go (source SHA-256: `2eec0a58e170f14b054a9d0c006fede8cd5c075c5d7cde1fbab688b4b793588e`)
- golang.org/x/sys@v0.49.0/cpu/cpu_x86.go (source SHA-256: `11df8d8a344644766603b3094ce22ff772a4f642e2dee7275ced0959a0083e38`)
- golang.org/x/sys@v0.49.0/unix/affinity_linux.go (source SHA-256: `b6c998fdf681f97006b0dae3ba0bf5d6d05a9f86fcc94f1518e37eab9dbc990a`)
- golang.org/x/sys@v0.49.0/unix/aliases.go (source SHA-256: `53e7eeba0503ad62ec18cdf2ca51a1785249a8646354439c854148dc57c06fb5`)
- golang.org/x/sys@v0.49.0/unix/ioctl_unsigned.go (source SHA-256: `991b152e32cde753f3b64d8ebc1d18e9078a3615006bc79c474b2bbbc653487a`)
- golang.org/x/sys@v0.49.0/unix/syscall_darwin_libSystem.go (source SHA-256: `0c327ad9b9845e19b1e097dfb7b569bc9793b670874e24be4a87ac9bd4647557`)
- golang.org/x/sys@v0.49.0/unix/syscall_linux_gc.go (source SHA-256: `41a54001d221205e9f1538066105b46c08cd329e85d4529426063961daa86235`)
- golang.org/x/sys@v0.49.0/windows/aliases.go (source SHA-256: `8209f09b0195574766435d50f0b769dd76a460777ff194c53952db22c4d54d77`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/cpu_gc_x86.s (source SHA-256: `74ac7fc7ef9c56c3306238cf031ea8ef7c0312a7116cb9fba6d07f0b1382df80`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/byteorder.go (source SHA-256: `825146fd4557b1cbd8161fa28bb4be8820089848d695316edeecb7fd5a551f8a`)
- golang.org/x/sys@v0.49.0/cpu/cpu_linux_noinit.go (source SHA-256: `620abd97199a2230f59165a2bc58c3ba780b6801ec0246cfb7319553dd795226`)
- golang.org/x/sys@v0.49.0/cpu/hwcap_linux.go (source SHA-256: `4101df793fddf76dfae477f917928008bc4a797446cd4ad44cde6e906c3f8713`)
- golang.org/x/sys@v0.49.0/unix/fcntl_darwin.go (source SHA-256: `5f690a6bac42da4bf5a688b9cec51de6fff9d9f60aebd9d9f7595e1f5fd5e16c`)
- golang.org/x/sys@v0.49.0/unix/fdset.go (source SHA-256: `8339dad44215930bdd253beaa1064c4d34ed2f81d54c4bd53d2be1f64fe35530`)
- golang.org/x/sys@v0.49.0/unix/readdirent_getdents.go (source SHA-256: `50e29bed47a256b79cb3754f43d0f8c1a5ea762e400bf809d0eb44d16ddd4714`)
- golang.org/x/sys@v0.49.0/unix/readdirent_getdirentries.go (source SHA-256: `db3a5b0e169d6f2431e3a145d2a3cc4374330e29d40378cb3fb26ede4fcd813d`)
- golang.org/x/sys@v0.49.0/unix/sockcmsg_unix_other.go (source SHA-256: `73f23d4f0c002e24160530515a45f46bb6fbe8728cc79b135a8c6956d5c8458c`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/ptrace_darwin.go (source SHA-256: `61eab9ab1e3d5d24b8e3545352e97afa0f76a77a578f75fb078026617444052c`)

````text
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/ifreq_linux.go (source SHA-256: `20ec91bd7bcd5e71f45316f23e7180bb6994aca7a87e49c253e4549e7436cbd4`)
- golang.org/x/sys@v0.49.0/unix/ioctl_linux.go (source SHA-256: `fd7dbee47ea5f00207bf340db54d294960101f003b6722d31162a06d69f4c4ff`)
- golang.org/x/sys@v0.49.0/unix/sysvshm_linux.go (source SHA-256: `8f932aac3521059ca3b4156c78b1ec45f1976972d1f70dee14375a1b7a507e44`)
- golang.org/x/sys@v0.49.0/unix/sysvshm_unix.go (source SHA-256: `e9e5031e048cf7c58692650ea14887158c3bedaaba506be63cbd8190cea2c066`)
- golang.org/x/sys@v0.49.0/unix/sysvshm_unix_other.go (source SHA-256: `b513c4e9cd077df2b1452bd645abcc3e252ec28ef44cf4e2a984d94bd5464fb3`)
- golang.org/x/sys@v0.49.0/windows/setupapi_windows.go (source SHA-256: `72097fc028210a5700ede92673d155af6de97210292770bfb31363691dad07c4`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/asm_bsd_amd64.s (source SHA-256: `1fecd01932d872c0d4ec06178a1860ae12bfab8490056dcf0a9d7a16cd455531`)
- golang.org/x/sys@v0.49.0/unix/asm_bsd_arm64.s (source SHA-256: `f7740a9d925eccd280e54e7971a36508a7d2856d9ef996a394ad5cfd80bec8c3`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/parse.go (source SHA-256: `97b269ea4f0b6d4071a9eed8a74f05055965c307bbab9090d9002bb01f7365a9`)
- golang.org/x/sys@v0.49.0/unix/syscall_linux_alarm.go (source SHA-256: `eede15c54c7c4dfc0691ae0e646fd31aca15bfd72f53783081e829cd150b19e0`)

````text
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/endian_little.go (source SHA-256: `c6bc70c372d9e1fe86fcf295f406b17bf04bf8d1af25c2456f58520cdaef3be9`)
- golang.org/x/sys@v0.49.0/cpu/runtime_auxv.go (source SHA-256: `d898ace395866bed261d403c8cd0ea6eab6d6d77f52042204957e389c38938cf`)
- golang.org/x/sys@v0.49.0/cpu/runtime_auxv_go121.go (source SHA-256: `6eee9d1a593dce53d22545dc5d5f6ed9127a43e6568ef5f5521920946510d445`)
- golang.org/x/sys@v0.49.0/unix/mmap_nomremap.go (source SHA-256: `f6166771a7d9a6c116617b6e21a15b020ddf709b39162c8596e0880cef3ce4be`)
- golang.org/x/sys@v0.49.0/unix/mremap.go (source SHA-256: `c4f762c0f38d6a9f08c9b754b207a1ba64dd29d03692245171f8a344aac6935a`)

````text
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/cpu_darwin_x86.go (source SHA-256: `4d90cd4b64eed4f9bbfe8e6c4b856f785659d6f22f56f8f6cf1b172c10e9c518`)
- golang.org/x/sys@v0.49.0/cpu/cpu_other_x86.go (source SHA-256: `4773a7cdfa25670c552bdd78cbe2bd874bd4b8f2d90f614c325ae79d2ffc0a7e`)
- golang.org/x/sys@v0.49.0/cpu/syscall_darwin_x86_gc.go (source SHA-256: `66059b45534544bdcde771fa793eca91d04de079e21adc4b2732e0912a57900a`)
- golang.org/x/sys@v0.49.0/unix/vgetrandom_linux.go (source SHA-256: `4dfa240c42977f70ef32d1cea7effdbd902b5ef6814e67f20c48da10548b4bc9`)
- golang.org/x/sys@v0.49.0/unix/vgetrandom_unsupported.go (source SHA-256: `822c28801c556c34f31cd7a14ddbe4c3708e543de6fe81c4ad1fae9420268258`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/asm_darwin_x86_gc.s (source SHA-256: `21022bf7461861b3744775075f4a4ef960c036cc469914684f87a4da9fc8b3a8`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/unix/auxv.go (source SHA-256: `5e470a481610ff746d64cb22b3e7a981ffa527d6ca546e87df4296704d6c6de6`)

````text
// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/sys@v0.49.0/cpu/cpu_windows.go (source SHA-256: `af3d72b6e5d5b0d6e801eaff989e8a38f6d26ad253faf057fa618682a1875a36`)
- golang.org/x/sys@v0.49.0/unix/readv_unix.go (source SHA-256: `2abd5ef6af8852208f2bae9de6de1c8d90ae33caba76f7fb757004d7a4040fe8`)

````text
// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## golang.org/x/term@v0.47.0

Source: https://pkg.go.dev/golang.org/x/term@v0.47.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:cpyOwQkr2Nzxe2p6KnrTG7Ei2wtzjsuSjhpXXepZ/BM=`

### golang.org/x/term@v0.47.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### golang.org/x/term@v0.47.0/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- golang.org/x/term@v0.47.0/terminal.go (source SHA-256: `221e717be6277af5e7cca327c9400577d81c6eb66c3606812b5d4a5dc58b8de7`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/term@v0.47.0/term_unix_bsd.go (source SHA-256: `582326355334ff9492a712458087a16b705c5855d44c30ed1ae45c0aeb0be7af`)

````text
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/term@v0.47.0/term.go (source SHA-256: `9792d91be6b3e5195591a2e5bda1fd01cb19954df79fd844de4d429538c1c5d5`)
- golang.org/x/term@v0.47.0/term_unix.go (source SHA-256: `e6a95d1ac948b5e2b389676545101de744f2715766e77f55160375a49a765b15`)
- golang.org/x/term@v0.47.0/term_windows.go (source SHA-256: `6a39deb1874975980a29e47c76b0bd87ed7f03d7955befaf478680a27deaa229`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/term@v0.47.0/term_unix_other.go (source SHA-256: `b85131ef51f97355fdadfec6f993c8fb4c0cd5c3d2453e1eec76c48d1d4d87d0`)

````text
// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## golang.org/x/text@v0.43.0

Source: https://pkg.go.dev/golang.org/x/text@v0.43.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:1QivrlhwAsnMOcqgOdl9my7s9OyZl8O5/UGUTJVXoWc=`

### golang.org/x/text@v0.43.0/LICENSE

SHA-256: `911f8f5782931320f5b8d1160a76365b83aea6447ee6c04fa6d5591467db9dad`

````text
Copyright 2009 The Go Authors.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google LLC nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### golang.org/x/text@v0.43.0/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- golang.org/x/text@v0.43.0/unicode/norm/composition.go (source SHA-256: `326e91fea709726d884e475db9b7a61ad17829e32d795f182cdcfcfdf64d9bf6`)
- golang.org/x/text@v0.43.0/unicode/norm/forminfo.go (source SHA-256: `b53ecff9ab5a83fd45598f54dfc9ea03849a3bcf4617ef02cf43225c71acc1e3`)
- golang.org/x/text@v0.43.0/unicode/norm/input.go (source SHA-256: `965b431790bb139543d71a9c497920ef7d9a15af417456a2bbd0cdb629330e8d`)
- golang.org/x/text@v0.43.0/unicode/norm/iter.go (source SHA-256: `4d580123776d78ff862131bd8c99fa5758f3cb70530619edc80ed4ee173cd83a`)
- golang.org/x/text@v0.43.0/unicode/norm/normalize.go (source SHA-256: `f1363700e71a35f9c966492f10dc6f3405cd824f2cc03b08ee64b24503d56f39`)
- golang.org/x/text@v0.43.0/unicode/norm/readwriter.go (source SHA-256: `d600d1f6cff2536dbb49a0f6bb8c06f6fc05f2754dad2dc05f1d44f9e8e48a35`)
- golang.org/x/text@v0.43.0/unicode/norm/trie.go (source SHA-256: `d87793d558251ee8824954f0b7bc5564803e4c9d59a8c4eebb4c3c5cfbd19492`)

````text
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/encoding/charmap/charmap.go (source SHA-256: `8f2418e924ccf719fcef0ec4ca1895657c1b98dcb7df4bfeba204970ebd98b33`)
- golang.org/x/text@v0.43.0/encoding/encoding.go (source SHA-256: `7a66d04774bf136fc1e655e0746eee54237e33ce21802d7e286c87b6e1307172`)
- golang.org/x/text@v0.43.0/encoding/japanese/eucjp.go (source SHA-256: `aa30243e36cc4dfa7168c31c7cce228109491f57b93f843d3832727c474fd50d`)
- golang.org/x/text@v0.43.0/encoding/japanese/iso2022jp.go (source SHA-256: `fa59e403d475ad0b185430c793fbd644167a2d6fa38b37670a18aa90b358c0ae`)
- golang.org/x/text@v0.43.0/encoding/japanese/shiftjis.go (source SHA-256: `23bf15129f276542aceefa460b68bd9905e68c460691b06e09c1b57c961694c3`)
- golang.org/x/text@v0.43.0/encoding/korean/euckr.go (source SHA-256: `9fbe514d8004ae60b1fc96bdb5212bbe9d06beffe50a2b0c8c20a83779ceabf9`)
- golang.org/x/text@v0.43.0/encoding/simplifiedchinese/gbk.go (source SHA-256: `509db1e531aade77b939704b307331c213b566122cb58fd6772ff16b0d4c2117`)
- golang.org/x/text@v0.43.0/encoding/simplifiedchinese/hzgb2312.go (source SHA-256: `4a7d1abc4120298617768ee17a70e0641b2c7149a86cdbb44189386fd19bd713`)
- golang.org/x/text@v0.43.0/encoding/traditionalchinese/big5.go (source SHA-256: `70eb462895050cd097282354ed49c944567ac14fe1552325c8eb0ace9b4c938b`)
- golang.org/x/text@v0.43.0/encoding/unicode/unicode.go (source SHA-256: `277d0a491016a630f27ee7d5ee6eae73abe6a431e0b3a3962e69e7dfd4f151ef`)
- golang.org/x/text@v0.43.0/internal/language/compact/language.go (source SHA-256: `1d49845ec8c93b000391718e3127c0979c054186570805fce27f7fe167ad2c5d`)
- golang.org/x/text@v0.43.0/internal/language/compact/tags.go (source SHA-256: `e6f98b11ee9c2d5d4d51a5d810d4ea2eff0872a8d11b962979744899e020e12c`)
- golang.org/x/text@v0.43.0/internal/language/language.go (source SHA-256: `19a0ccf7635fd3dbad3766d953996df5d1c24a367e44c49200ab6ec3023fc871`)
- golang.org/x/text@v0.43.0/internal/language/lookup.go (source SHA-256: `8f0ac8780f89d5c2b29b98fd599269b445d7bb612a7948c9e97d425424b5ea38`)
- golang.org/x/text@v0.43.0/internal/language/match.go (source SHA-256: `711e9c548516929762f0636a69788ba6eed6ed10bfb47634f3b5d5a1f26a1547`)
- golang.org/x/text@v0.43.0/internal/language/parse.go (source SHA-256: `0c884b51bb83d3633a51f16e90c7dbaf45c301a449e366df591eda11a3000972`)
- golang.org/x/text@v0.43.0/internal/language/tags.go (source SHA-256: `21041c1c1982b12db6f523a4feb58746dac3e2c71e8dcf9f76ac5807cd809099`)
- golang.org/x/text@v0.43.0/language/language.go (source SHA-256: `3f2032f0367d83fb14772384d47119d27c2c427b277aa0a02913db59abe32bb0`)
- golang.org/x/text@v0.43.0/language/match.go (source SHA-256: `b8a6b973391cba46899d079173cf0b1d3064f28bb79a11c3fdc65c089cbe23d3`)
- golang.org/x/text@v0.43.0/language/parse.go (source SHA-256: `3552f8fba45f310216182d167ba1da8e0c02212ad1896eeb66b194dd546b7bc2`)
- golang.org/x/text@v0.43.0/language/tags.go (source SHA-256: `247c5ce419821c135a2751029c078cc6e1472c9896dd581d67a277449e54d968`)
- golang.org/x/text@v0.43.0/transform/transform.go (source SHA-256: `7716301e210f42b764ec78f4adffb05075f1d5103446a40b8a38043f5a7abd88`)
- golang.org/x/text@v0.43.0/unicode/norm/transform.go (source SHA-256: `59aa584156c76ddb33f45d8d2c54818f3d42287a382b2fcffebec82085c9bc0d`)

````text
// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/internal/language/coverage.go (source SHA-256: `8c8238012deb1b3c42820c755343db1bc326dc5cbddfcc24f97abbb7e0c6cee1`)
- golang.org/x/text@v0.43.0/language/coverage.go (source SHA-256: `d12716948854471a3536a04dc48736b0db0318ea6df61db236bc4ee0fe725375`)
- golang.org/x/text@v0.43.0/runes/runes.go (source SHA-256: `30879876f51d507e1cce28de3aa6c63966d101f92df42e6e0be18fec9dd7592f`)

````text
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/encoding/htmlindex/htmlindex.go (source SHA-256: `f7de101bd433f74cc2d32a92d224c36cf04225b48323c4d66ab79805a59fb6c0`)
- golang.org/x/text@v0.43.0/encoding/htmlindex/map.go (source SHA-256: `8094e6c1a15cc74160ecd3ed326d38fdc55af9e7b1caef53f5370e31f1e94cd5`)
- golang.org/x/text@v0.43.0/encoding/internal/identifier/identifier.go (source SHA-256: `50c99f12c73517c9314d04fb7807e7fdc57e75436306d3b38036fb17bc0c41ca`)
- golang.org/x/text@v0.43.0/encoding/internal/internal.go (source SHA-256: `d49433c5ac6f01c29013957898b7fd9d899832d2c7d09cc137ff52a4270aff87`)
- golang.org/x/text@v0.43.0/encoding/japanese/all.go (source SHA-256: `5728aaadf3e9b98632c5cbbedc60e29ec14779bc8f020bae135c97b1548e3944`)
- golang.org/x/text@v0.43.0/encoding/simplifiedchinese/all.go (source SHA-256: `98b67e8457b51b3b94fb0a1276c3aa006ce4bebe91f4911787d44bf964ae7898`)
- golang.org/x/text@v0.43.0/encoding/unicode/override.go (source SHA-256: `56800213c1d96ac93e091ac0bdf084caf259fac616fcea62b49c5f5ecff24b8f`)
- golang.org/x/text@v0.43.0/internal/tag/tag.go (source SHA-256: `97a27a6bfdd26d872e07d051fc27e622e1cd764943350e8bfa9cf0e2e13cc474`)
- golang.org/x/text@v0.43.0/internal/utf8internal/utf8internal.go (source SHA-256: `a2aa5c3b12a72996fcedd945dd7d1955f2e5ef76bde6f5277238d2b2d6506390`)
- golang.org/x/text@v0.43.0/runes/cond.go (source SHA-256: `a97ca46dc0d82c3f04b17efdb26244f9f779ef1cabcdd5875b17fa51dc24640e`)
- golang.org/x/text@v0.43.0/unicode/bidi/bidi.go (source SHA-256: `b1e9360e9a026006f56914a5716b18780346418bb0b64bce23a4ff92d8147b13`)
- golang.org/x/text@v0.43.0/unicode/bidi/bracket.go (source SHA-256: `5caa24f9c4e04a54d18875468f976b4b7c09c8700e88286f2fd96828f134107e`)
- golang.org/x/text@v0.43.0/unicode/bidi/core.go (source SHA-256: `ba8439de7125be9ec4c47439ce6c517f73277d32313cd1c1b723bc8deee88aca`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/secure/bidirule/bidirule.go (source SHA-256: `6f6dd76ab63af3ed24759fd39328026a3f40c9739e45cf152b2f02a420f4e310`)
- golang.org/x/text@v0.43.0/unicode/bidi/prop.go (source SHA-256: `f6391b2f69a1ae2a0ac1f0599a90b260ec350a19f73ad5d6572660100b5950bc`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/language/doc.go (source SHA-256: `4b97559b103832f3c5f22e80f2df0dd7394a469d945f79afad1654eb179fa63b`)

````text
// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/internal/language/compact.go (source SHA-256: `2d121c703729c5fbd5f3f13fd568af142ca8fde0f01a07e1874df22dbaadde98`)
- golang.org/x/text@v0.43.0/internal/language/compact/compact.go (source SHA-256: `e154aa511d328fd47a9c9fc1115d9c92be2d385e885ae65c5e348ec1d5844ead`)
- golang.org/x/text@v0.43.0/internal/language/compose.go (source SHA-256: `7657b95eac8448d5d2273c048d316dde3ac0e49ac058487a69b2b990ebc7f8ef`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- golang.org/x/text@v0.43.0/encoding/internal/identifier/mib.go (source SHA-256: `7dac7c2f13c9a8f31f6f8228e71ea0ae6f470f3099a2ba8cebc4a1fbc876385a`)

````text
// VenturaUS is the MIB identifier with IANA name Ventura-US.
	//
	// Ventura US.  ASCII plus characters typically used in
	// publishing, like pilcrow, copyright, registered, trade mark,
	// section, dagger, and double dagger in the range A0 (hex)
	// to FF (hex).
	// PCL Symbol Set id: 14J
````


## google.golang.org/api@v0.308.0

Source: https://pkg.go.dev/google.golang.org/api@v0.308.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:iGb+4qT3nVb3uuIzwor3yP2nyu+rMZr31PvLFiRNJQk=`

### google.golang.org/api@v0.308.0/AUTHORS

SHA-256: `7c195dee59dd2338b87926657674cc39d06c7c93b8ae13c86ca07749a4be50ad`

````text
# This is the official list of authors for copyright purposes.
# This file is distinct from the CONTRIBUTORS files.
# See the latter for an explanation.

# Names should be added to this file as
#	Name or Organization <email address>
# The email address is not required for organizations.

# Please keep the list sorted.
Google Inc.
LightStep Inc.
````

### google.golang.org/api@v0.308.0/LICENSE

SHA-256: `110244b02140866ee37d17fa7449436a377ec3b85a481fbb208f4c87964382de`

````text
Copyright (c) 2011 Google Inc. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### google.golang.org/api@v0.308.0/internal/third_party/uritemplates/LICENSE

SHA-256: `fc0a2f71df4e8f047902da53d1f85301be43e0f360fc167057a2d04658ed2ba9`

````text
Copyright (c) 2013 Joshua Tacoma. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### Original source notice

- google.golang.org/api@v0.308.0/googleapi/googleapi.go (source SHA-256: `f5e3aecc7a1c302cca7d78aa7a52468521d7682edb991279ba5798621357d7ae`)

````text
// Copyright 2011 Google LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/googleapi/transport/apikey.go (source SHA-256: `359c95027c6ac83d8820d7f76832952357c18ad2c3059410421ea1260f7ad50c`)

````text
// Copyright 2012 Google LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/googleapi/types.go (source SHA-256: `5d7e91863dc82318435e36d9e2ef16eaed75c4323abbe25f7c66aa58c7780435`)

````text
// Copyright 2013 Google LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/third_party/uritemplates/uritemplates.go (source SHA-256: `13be84b06c681656a3af14210f387af44f7dea1c4110276471c14e9245a74a5e`)

````text
// Copyright 2013 Joshua Tacoma. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/transport/http/dial.go (source SHA-256: `ed347046dd0849a756eb9f331ea094d7aab5109a3b7d8e595b3c523845b96265`)

````text
// Copyright 2015 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/gensupport/json.go (source SHA-256: `9176cf4ec671813864de1db0206926cef8cbe30689597d60f9273963669563b0`)
- google.golang.org/api@v0.308.0/internal/gensupport/params.go (source SHA-256: `d05bca2a4cc8eb59548688f278d71e8d2bedec9200c05fa057408108d2a3c352`)

````text
// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/gensupport/jsonfloat.go (source SHA-256: `a44e90e16c501da55ea7e8d912ff254ae5e7b314947d2e0df7cdf5eaf65677b8`)

````text
// Copyright 2016 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/gensupport/buffer.go (source SHA-256: `af068ff72cd4c273e5a670bc3fdf85c9af1e374f2866c1cf9164d8203c921768`)
- google.golang.org/api@v0.308.0/internal/gensupport/doc.go (source SHA-256: `fb0b7bc8969a371a770b656d2e095b2a988052463dbeabaf0dc0383d38f880be`)
- google.golang.org/api@v0.308.0/internal/gensupport/media.go (source SHA-256: `47c81438de78d304c1c2076dddb6a745ac01fe63b6193feee3cb8f8a2bc1003d`)
- google.golang.org/api@v0.308.0/internal/gensupport/resumable.go (source SHA-256: `f02844474c0dd9aef774c2f6605f93b58caa87845f786c37633cdd543fb43141`)
- google.golang.org/api@v0.308.0/internal/gensupport/send.go (source SHA-256: `fdaacb3e456e80dafc5a49fcda26b252f113e7deece669c009b1ae28e2cee7a5`)
- google.golang.org/api@v0.308.0/internal/third_party/uritemplates/utils.go (source SHA-256: `f1552b96f8d4356df20757f8c9821dab24fe8e4cc0d1544fb6330911788a0313`)

````text
// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/creds.go (source SHA-256: `7067ca2f811fbfba5758ee20cc4a778fc0acedf553028d9c42f0410d88c3fe55`)
- google.golang.org/api@v0.308.0/internal/settings.go (source SHA-256: `a3a3489bd987755d2cfbdc8ff7f1a68e513ea423d50b4913813f9c823eaa3da6`)
- google.golang.org/api@v0.308.0/option/option.go (source SHA-256: `f9ae99d5087a99a630522c257541a864c3bce75f3acec0c932661c77ab68c0b3`)

````text
// Copyright 2017 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/cba.go (source SHA-256: `0f6443fab04eead1a08ea81efb6cd265ecd6835bf1831f32fcfae644a51f0bcb`)
- google.golang.org/api@v0.308.0/internal/cert/default_cert.go (source SHA-256: `e4df83b81997cc8516d71eee25230074f893f6e5e15159b9e155cab2e8e282fe`)
- google.golang.org/api@v0.308.0/internal/conn_pool.go (source SHA-256: `d064443f748e1dda02b6a22325c84c13eb0118803d11c1eff830dac6b383cb96`)
- google.golang.org/api@v0.308.0/option/internaloption/internaloption.go (source SHA-256: `9de42b258ee7b6cfbcf60ac7698d8cb71d199fdecdd21153da7acb0f50efb234`)

````text
// Copyright 2020 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/gensupport/version.go (source SHA-256: `5ad1debeff2a611fc3d1508775c920a573a40b19cb06d5c7b8017d40d7fa45df`)

````text
// Copyright 2020 Google LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/impersonate/impersonate.go (source SHA-256: `8c63a0034b4c3136acfaa2d67ea9c562667180cdda285151084fd7be052cc8e0`)

````text
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/gensupport/retry.go (source SHA-256: `5cd6d3c60891aedb1843bf95a7cb7271804811c9cd09555e1784e3b1414df062`)

````text
// Copyright 2021 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/cert/enterprise_cert.go (source SHA-256: `6ac02e3004896b147a11b7d673e4fe7fd36d104be15c9f0aeed816795c91c205`)
- google.golang.org/api@v0.308.0/internal/cert/secureconnect_cert.go (source SHA-256: `c2a0df0f30dadf64e5151ce59e940e0ece17fd67468043b699f07d236cc4d46c`)

````text
// Copyright 2022 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/gensupport/error.go (source SHA-256: `8537e8e823949559a055c4d6f19f635e91df8d72f92a7251648c4537ed93cfb2`)
- google.golang.org/api@v0.308.0/internal/version.go (source SHA-256: `d5b51cf45c8f1a437cd4061cb6d13cb8c28f6e31d4bdea15623a2e8e57a09194`)

````text
// Copyright 2022 Google LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/s2a.go (source SHA-256: `b329821a611d38fd04961ea80f927173ed17d00eaa365a2997c3f2fd7a3cd904`)

````text
// Copyright 2023 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/internal/credentialstype/credentialstype.go (source SHA-256: `550dc2f0e17b7d1bed4ab3cb24670d035ed92a13d9bd9c7d54cd168914f0ac54`)

````text
// Copyright 2024 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/api@v0.308.0/dataportability/v1/dataportability-gen.go (source SHA-256: `ed104cc6b6651da2153f3d7fbf1257dcaaaec81f9b0d87d46f07aec1afd380c6`)
- google.golang.org/api@v0.308.0/drive/v3/drive-gen.go (source SHA-256: `4b52810b95c0d9e4b0f5e3ba01f1d826d0857a3b98366f5e1615dd29c2ddefbd`)
- google.golang.org/api@v0.308.0/option/internaloption/unsaferesolver.go (source SHA-256: `939c1de0e24925b161af555221b0f51c85caefcdd4443419dd94fefab23c7055`)

````text
// Copyright 2026 Google LLC.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````


## google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b

Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:6WNWHZk91dp4gZ25642OALkuG1UVNZbwENJD19DnB8k=`

### google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b/LICENSE

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### Original source notice

- google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b/code/code.pb.go (source SHA-256: `c2350de979dcce2072fecbb9aa1e1d717dac660fcab04086a841a728b9b79c7f`)
- google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b/errdetails/error_details.pb.go (source SHA-256: `8034b9c7462d9af0eb10e522841a53e508517145611ef366efa1bd4bf4fd74ae`)
- google.golang.org/genproto/googleapis/rpc@v0.0.0-20261019180522-ce908d29299b/status/status.pb.go (source SHA-256: `2f3ea7d76f05c49d39ab22bae8d99cf45b847fd2a99cdc2a4637a3857a4fdbf0`)

````text
// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## google.golang.org/grpc@v1.84.0

Source: https://pkg.go.dev/google.golang.org/grpc@v1.84.0

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:soMyaPJ8pAak5PIQ0DGBUir0XRo2fRoMqhNWMLlLxO0=`

### google.golang.org/grpc@v1.84.0/AUTHORS

SHA-256: `627c695f80fdac6412d66e6ad4160a4acb41923ca7ee3fd7e112e105f0ff46b9`

````text
Google Inc.
````

### google.golang.org/grpc@v1.84.0/LICENSE

SHA-256: `cfc7749b96f63bd31c3c42b5c471bf756814053e847c10f3eb003417bc523d30`

````text

                                 Apache License
                           Version 2.0, January 2004
                        http://www.apache.org/licenses/

   TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION

   1. Definitions.

      "License" shall mean the terms and conditions for use, reproduction,
      and distribution as defined by Sections 1 through 9 of this document.

      "Licensor" shall mean the copyright owner or entity authorized by
      the copyright owner that is granting the License.

      "Legal Entity" shall mean the union of the acting entity and all
      other entities that control, are controlled by, or are under common
      control with that entity. For the purposes of this definition,
      "control" means (i) the power, direct or indirect, to cause the
      direction or management of such entity, whether by contract or
      otherwise, or (ii) ownership of fifty percent (50%) or more of the
      outstanding shares, or (iii) beneficial ownership of such entity.

      "You" (or "Your") shall mean an individual or Legal Entity
      exercising permissions granted by this License.

      "Source" form shall mean the preferred form for making modifications,
      including but not limited to software source code, documentation
      source, and configuration files.

      "Object" form shall mean any form resulting from mechanical
      transformation or translation of a Source form, including but
      not limited to compiled object code, generated documentation,
      and conversions to other media types.

      "Work" shall mean the work of authorship, whether in Source or
      Object form, made available under the License, as indicated by a
      copyright notice that is included in or attached to the work
      (an example is provided in the Appendix below).

      "Derivative Works" shall mean any work, whether in Source or Object
      form, that is based on (or derived from) the Work and for which the
      editorial revisions, annotations, elaborations, or other modifications
      represent, as a whole, an original work of authorship. For the purposes
      of this License, Derivative Works shall not include works that remain
      separable from, or merely link (or bind by name) to the interfaces of,
      the Work and Derivative Works thereof.

      "Contribution" shall mean any work of authorship, including
      the original version of the Work and any modifications or additions
      to that Work or Derivative Works thereof, that is intentionally
      submitted to Licensor for inclusion in the Work by the copyright owner
      or by an individual or Legal Entity authorized to submit on behalf of
      the copyright owner. For the purposes of this definition, "submitted"
      means any form of electronic, verbal, or written communication sent
      to the Licensor or its representatives, including but not limited to
      communication on electronic mailing lists, source code control systems,
      and issue tracking systems that are managed by, or on behalf of, the
      Licensor for the purpose of discussing and improving the Work, but
      excluding communication that is conspicuously marked or otherwise
      designated in writing by the copyright owner as "Not a Contribution."

      "Contributor" shall mean Licensor and any individual or Legal Entity
      on behalf of whom a Contribution has been received by Licensor and
      subsequently incorporated within the Work.

   2. Grant of Copyright License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      copyright license to reproduce, prepare Derivative Works of,
      publicly display, publicly perform, sublicense, and distribute the
      Work and such Derivative Works in Source or Object form.

   3. Grant of Patent License. Subject to the terms and conditions of
      this License, each Contributor hereby grants to You a perpetual,
      worldwide, non-exclusive, no-charge, royalty-free, irrevocable
      (except as stated in this section) patent license to make, have made,
      use, offer to sell, sell, import, and otherwise transfer the Work,
      where such license applies only to those patent claims licensable
      by such Contributor that are necessarily infringed by their
      Contribution(s) alone or by combination of their Contribution(s)
      with the Work to which such Contribution(s) was submitted. If You
      institute patent litigation against any entity (including a
      cross-claim or counterclaim in a lawsuit) alleging that the Work
      or a Contribution incorporated within the Work constitutes direct
      or contributory patent infringement, then any patent licenses
      granted to You under this License for that Work shall terminate
      as of the date such litigation is filed.

   4. Redistribution. You may reproduce and distribute copies of the
      Work or Derivative Works thereof in any medium, with or without
      modifications, and in Source or Object form, provided that You
      meet the following conditions:

      (a) You must give any other recipients of the Work or
          Derivative Works a copy of this License; and

      (b) You must cause any modified files to carry prominent notices
          stating that You changed the files; and

      (c) You must retain, in the Source form of any Derivative Works
          that You distribute, all copyright, patent, trademark, and
          attribution notices from the Source form of the Work,
          excluding those notices that do not pertain to any part of
          the Derivative Works; and

      (d) If the Work includes a "NOTICE" text file as part of its
          distribution, then any Derivative Works that You distribute must
          include a readable copy of the attribution notices contained
          within such NOTICE file, excluding those notices that do not
          pertain to any part of the Derivative Works, in at least one
          of the following places: within a NOTICE text file distributed
          as part of the Derivative Works; within the Source form or
          documentation, if provided along with the Derivative Works; or,
          within a display generated by the Derivative Works, if and
          wherever such third-party notices normally appear. The contents
          of the NOTICE file are for informational purposes only and
          do not modify the License. You may add Your own attribution
          notices within Derivative Works that You distribute, alongside
          or as an addendum to the NOTICE text from the Work, provided
          that such additional attribution notices cannot be construed
          as modifying the License.

      You may add Your own copyright statement to Your modifications and
      may provide additional or different license terms and conditions
      for use, reproduction, or distribution of Your modifications, or
      for any such Derivative Works as a whole, provided Your use,
      reproduction, and distribution of the Work otherwise complies with
      the conditions stated in this License.

   5. Submission of Contributions. Unless You explicitly state otherwise,
      any Contribution intentionally submitted for inclusion in the Work
      by You to the Licensor shall be under the terms and conditions of
      this License, without any additional terms or conditions.
      Notwithstanding the above, nothing herein shall supersede or modify
      the terms of any separate license agreement you may have executed
      with Licensor regarding such Contributions.

   6. Trademarks. This License does not grant permission to use the trade
      names, trademarks, service marks, or product names of the Licensor,
      except as required for reasonable and customary use in describing the
      origin of the Work and reproducing the content of the NOTICE file.

   7. Disclaimer of Warranty. Unless required by applicable law or
      agreed to in writing, Licensor provides the Work (and each
      Contributor provides its Contributions) on an "AS IS" BASIS,
      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
      implied, including, without limitation, any warranties or conditions
      of TITLE, NON-INFRINGEMENT, MERCHANTABILITY, or FITNESS FOR A
      PARTICULAR PURPOSE. You are solely responsible for determining the
      appropriateness of using or redistributing the Work and assume any
      risks associated with Your exercise of permissions under this License.

   8. Limitation of Liability. In no event and under no legal theory,
      whether in tort (including negligence), contract, or otherwise,
      unless required by applicable law (such as deliberate and grossly
      negligent acts) or agreed to in writing, shall any Contributor be
      liable to You for damages, including any direct, indirect, special,
      incidental, or consequential damages of any character arising as a
      result of this License or out of the use or inability to use the
      Work (including but not limited to damages for loss of goodwill,
      work stoppage, computer failure or malfunction, or any and all
      other commercial damages or losses), even if such Contributor
      has been advised of the possibility of such damages.

   9. Accepting Warranty or Additional Liability. While redistributing
      the Work or Derivative Works thereof, You may choose to offer,
      and charge a fee for, acceptance of support, warranty, indemnity,
      or other liability obligations and/or rights consistent with this
      License. However, in accepting such obligations, You may act only
      on Your own behalf and on Your sole responsibility, not on behalf
      of any other Contributor, and only if You agree to indemnify,
      defend, and hold each Contributor harmless for any liability
      incurred by, or claims asserted against, such Contributor by reason
      of your accepting any such warranty or additional liability.

   END OF TERMS AND CONDITIONS

   APPENDIX: How to apply the Apache License to your work.

      To apply the Apache License to your work, attach the following
      boilerplate notice, with the fields enclosed by brackets "[]"
      replaced with your own identifying information. (Don't include
      the brackets!)  The text should be enclosed in the appropriate
      comment syntax for the file format. We also recommend that a
      file or class name and description of purpose be included on the
      same "printed page" as the copyright notice for easier
      identification within third-party archives.

   Copyright [yyyy] [name of copyright owner]

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
````

### google.golang.org/grpc@v1.84.0/NOTICE.txt

SHA-256: `693ff28ec216d5112ac1bbfe64ef539005867d1c7bd427b57d579683293b947f`

````text
Copyright 2014 gRPC authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
````

### Original source notice

- google.golang.org/grpc@v1.84.0/call.go (source SHA-256: `9b156edda59b97671c0576ce9d0a283d7afb8dc24cc621b4a5b9089bc7f49f2e`)
- google.golang.org/grpc@v1.84.0/clientconn.go (source SHA-256: `8182ba3e972eaf9e93f91c07969b2c795f8bd8942d99e8f66413f599c5883a0f`)
- google.golang.org/grpc@v1.84.0/codec.go (source SHA-256: `6717f2e6bd13fd93c6f6726f66153d0ea132a8aec7c45c6258b7375727055955`)
- google.golang.org/grpc@v1.84.0/codes/codes.go (source SHA-256: `e67aa5a4a8a39663b70244ee9ef6f76b8ef145896cdf3cd10b397d9173aa7322`)
- google.golang.org/grpc@v1.84.0/credentials/credentials.go (source SHA-256: `eacbcc8137ade7617e88ebdab497db2f4b3b0cb4d87543ee0963098120ec806b`)
- google.golang.org/grpc@v1.84.0/credentials/tls.go (source SHA-256: `8f08b7d81dbef01c9349540bfd0ef43775c12bcba79a6319d20f14d190e4022f`)
- google.golang.org/grpc@v1.84.0/internal/transport/controlbuf.go (source SHA-256: `48c62d0504d7e5a1c7a2bbb05eab93076a4176b0c3267cc755d6a2db4cf59035`)
- google.golang.org/grpc@v1.84.0/internal/transport/flowcontrol.go (source SHA-256: `aff3b1bb108e20a280a2e2ae58ef77025ccffe2c4d60dc2bff843d03cf4e53c6`)
- google.golang.org/grpc@v1.84.0/internal/transport/http2_client.go (source SHA-256: `93fc9536fa02161c8559c21da7e26359d06d7ae4368a4ee7a4a6661cfe094491`)
- google.golang.org/grpc@v1.84.0/internal/transport/http2_server.go (source SHA-256: `6ff2da17e276ba22a782dffe467ec32b664895d12faccfbb5fc1ec40d59dcf6c`)
- google.golang.org/grpc@v1.84.0/internal/transport/http_util.go (source SHA-256: `b36a1b262f73a8eb45c9a3c3dae95710f2b4537673e15af6aa97e5e1c7f90545`)
- google.golang.org/grpc@v1.84.0/internal/transport/transport.go (source SHA-256: `40430d67f653b170ecb4fa9e73b94e2dbd697ac7251e8d7b934111a5163c89b4`)
- google.golang.org/grpc@v1.84.0/metadata/metadata.go (source SHA-256: `a9fff99819d67662d5c57dc3f2af600907ce3b8c3a95ead5feab12137d62512c`)
- google.golang.org/grpc@v1.84.0/peer/peer.go (source SHA-256: `9c62fce90f7e79f5e0f1e630aa207e3a1351fcd7c473710222df2acfd2f5bc51`)
- google.golang.org/grpc@v1.84.0/rpc_util.go (source SHA-256: `c8e8d24f53a44158529dc232db32ace4bc0f9ce67b9386e4df19abe935cded7e`)
- google.golang.org/grpc@v1.84.0/server.go (source SHA-256: `c920c82c4b7e8e28b00c799ff050cf42cbf03293984f2a65f04ed4323b30d273`)
- google.golang.org/grpc@v1.84.0/stream.go (source SHA-256: `68aa9293b8023d550c730d391482dc08b2572526b8629904f1dde457bf042c4d`)

````text
/*
 *
 * Copyright 2014 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/doc.go (source SHA-256: `cb8655c4e8c7342a9d06570b4ffe4706a42c0f7fa598c339757e2a6c2a7708d7`)
- google.golang.org/grpc@v1.84.0/grpclog/logger.go (source SHA-256: `2b59f8d4c49262754f1e73afd17763195f107fe16f00860083b63031a8ef3526`)
- google.golang.org/grpc@v1.84.0/trace.go (source SHA-256: `3c270cce29e756202e93dbdf6f60ccb4f3793b8657842cc8871f1791da7580cd`)

````text
/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/interceptor.go (source SHA-256: `75d0dab7b224efcbbd38e3fb2eb6724237c361357aa6948b0a857d535bf3bdeb`)
- google.golang.org/grpc@v1.84.0/internal/transport/handler_server.go (source SHA-256: `65c8274e186e499ad61541934983eeb5d43c29cebb2151cdae458d9b8cb4e73a`)
- google.golang.org/grpc@v1.84.0/stats/handlers.go (source SHA-256: `92dc01d539497278d78a81ad0edb631beda3fb15f4da81ff90b17068898ef057`)
- google.golang.org/grpc@v1.84.0/stats/stats.go (source SHA-256: `38d83e728848e3e24cd205bfd7e3319f96e5b300f793db1573cf58e49165e029`)
- google.golang.org/grpc@v1.84.0/tap/tap.go (source SHA-256: `e33fe01886824c60ac943591eb5da6aa09d46709c52a74ea9d4ba90d6ec3c6d9`)

````text
/*
 *
 * Copyright 2016 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/backoff.go (source SHA-256: `b5593815ced2bc309a68b9ab18c9044b6c93723ae303ef636f77c1ec7f400718`)
- google.golang.org/grpc@v1.84.0/balancer/balancer.go (source SHA-256: `189ebc1bb54397dafb3c2b7fcd6d2c5a60fc1e87f035773113b68c8b5be7a0ad`)
- google.golang.org/grpc@v1.84.0/balancer/base/balancer.go (source SHA-256: `a054d3ff2f1de786456e08f976a62de79036a2ecb44e9f32c58ff8b390b5f235`)
- google.golang.org/grpc@v1.84.0/balancer/base/base.go (source SHA-256: `e222da7aeb4c7341c59bcff314ec7c01b8900d3b2c27b3f1cb6dc5b55839f266`)
- google.golang.org/grpc@v1.84.0/balancer/pickfirst/pickfirst.go (source SHA-256: `df605f10bde58cc1a0e9b6c88cd68f6363a8c23463b827fb39ecb3531440bc3e`)
- google.golang.org/grpc@v1.84.0/balancer/roundrobin/roundrobin.go (source SHA-256: `2c23443fdc43aafe86d530124f6b5e52d65af86de5954265f33ec01a763d9c63`)
- google.golang.org/grpc@v1.84.0/balancer_wrapper.go (source SHA-256: `735c36b5c6fdb7a5f114615e4e49d02947458ed24a98f524b6e7917020371714`)
- google.golang.org/grpc@v1.84.0/codes/code_string.go (source SHA-256: `fc83f20f33a04a641ef92e72f2df310ad89875fed491ce402029308dcca53508`)
- google.golang.org/grpc@v1.84.0/connectivity/connectivity.go (source SHA-256: `ca86fd043275378e7c04ac49924a4307dd34f1fbef2a90b322af92cd0ee06d2e`)
- google.golang.org/grpc@v1.84.0/encoding/encoding.go (source SHA-256: `cf9661e3691d7bc0b52cd1263a6ba8c39f980869b6e946e3fc7313ce2284b0cb`)
- google.golang.org/grpc@v1.84.0/grpclog/grpclog.go (source SHA-256: `6b1c4c32181da0292783b4feb489b979830556724bfa3b33e4612b1279d24081`)
- google.golang.org/grpc@v1.84.0/grpclog/loggerv2.go (source SHA-256: `e5a0d19469e8fb4f14b809f118dbdbf14b76bfb0698679eac53f45d3d3793779`)
- google.golang.org/grpc@v1.84.0/internal/backoff/backoff.go (source SHA-256: `f6385bcafb72a06a86d5c28e0cd1680fdad8e3757182e2291fe13b495615a991`)
- google.golang.org/grpc@v1.84.0/internal/resolver/passthrough/passthrough.go (source SHA-256: `9b653000c8516ed008ed2416cbc59e48ffdf3b974330dae9e142004431392959`)
- google.golang.org/grpc@v1.84.0/internal/transport/bdp_estimator.go (source SHA-256: `01a7332c6ea43122862b91c7e69ca4a2aa91b11ee9c094fc2ffe1ba6edad75f3`)
- google.golang.org/grpc@v1.84.0/internal/transport/proxy.go (source SHA-256: `af52f2bac21e4f70ec01fae9bad77391eb4d1b5b4f4307f1a963bf33bfa7b2a2`)
- google.golang.org/grpc@v1.84.0/keepalive/keepalive.go (source SHA-256: `7ac071d59e43ba03b15fa318ed4c635be4c609a3c0091f8dde3e29573a29d7cc`)
- google.golang.org/grpc@v1.84.0/picker_wrapper.go (source SHA-256: `705305df53559ad0f775d081d262788ea6ed721500343730034ca980f8d56f5c`)
- google.golang.org/grpc@v1.84.0/resolver/resolver.go (source SHA-256: `19f7efcd95b6dc1d2d49bdf19ccd5b4ef4393b5c8f9f37133166ae7fd4a2bcc8`)
- google.golang.org/grpc@v1.84.0/resolver_wrapper.go (source SHA-256: `ba9547eed70c97448b4e69ffd2e48abcf333b92304cc59f788d824d364fa173b`)
- google.golang.org/grpc@v1.84.0/service_config.go (source SHA-256: `c38647649ab124f58b5e2c27dc7683afe6d6a277d3fc4966abe89c19e0299604`)
- google.golang.org/grpc@v1.84.0/status/status.go (source SHA-256: `bc5f3bd9fefa5b9c59c811e5b67bbd0351b9a042ed7013d6ecfb967511ff7b1c`)

````text
/*
 *
 * Copyright 2017 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/dialoptions.go (source SHA-256: `d94f52d64205f8e114ba3ac095ee497b98314d027564807239ce9c8d374c2ee0`)
- google.golang.org/grpc@v1.84.0/internal/binarylog/binarylog.go (source SHA-256: `5fc48a515beffa63f5e1ab297c8fd98ee77e755881fb392e55e313bd08191189`)
- google.golang.org/grpc@v1.84.0/internal/binarylog/binarylog_testutil.go (source SHA-256: `0490bfd3780ebf38e67f9b2c297759708f2f67b91a6db40b27baefa85efc258e`)
- google.golang.org/grpc@v1.84.0/internal/binarylog/env_config.go (source SHA-256: `960fd8f1f3159361bbae1110e4b9acd5d7d411534bd54b75723f85c64e176087`)
- google.golang.org/grpc@v1.84.0/internal/binarylog/method_logger.go (source SHA-256: `a4f12d36d3e6b89ca92bd77980936c1daad66031770bc814660d5bd503f2c970`)
- google.golang.org/grpc@v1.84.0/internal/binarylog/sink.go (source SHA-256: `7b76e9a8456f515cdc7a906157ea5b91fedb153239536abe8be528951b2b0b85`)
- google.golang.org/grpc@v1.84.0/internal/channelz/channelmap.go (source SHA-256: `d26c23599a3008b34dd1f94bda99cb4589d9e8a6a2c789c477d317df27f5e563`)
- google.golang.org/grpc@v1.84.0/internal/channelz/funcs.go (source SHA-256: `4aecad9253e8832541a4ee49136fcd4f1c0ba881972099d87caca20574537b1f`)
- google.golang.org/grpc@v1.84.0/internal/channelz/syscall_linux.go (source SHA-256: `6b90796f288344124b635c90237961789ca452ddf78e4b7f9c0417d0b458f61c`)
- google.golang.org/grpc@v1.84.0/internal/channelz/syscall_nonlinux.go (source SHA-256: `1be23873d86f82b13b2c7d8d2534b5380215120d15ce3e47bf8d25c9fbdf983c`)
- google.golang.org/grpc@v1.84.0/internal/channelz/trace.go (source SHA-256: `181af874a10e982864274e5a7476dbe5a35f9cc629e37820094c01fa99ee3dde`)
- google.golang.org/grpc@v1.84.0/internal/credentials/syscallconn.go (source SHA-256: `47cb93c2b159a3d1e9373a49a44dc50744e51f40ff3f16414cd15c7ceb390ef7`)
- google.golang.org/grpc@v1.84.0/internal/envconfig/envconfig.go (source SHA-256: `a989d4ee800ee59b386b3d01c4098b6dcc697894690d4a2b3168ac38e456614f`)
- google.golang.org/grpc@v1.84.0/internal/grpcsync/event.go (source SHA-256: `ff64826bdb0ec40990d49c72eec9903326f490b854900687303472964b00d06e`)
- google.golang.org/grpc@v1.84.0/internal/grpcutil/method.go (source SHA-256: `75cc6e1903494f229ef7feee6731cf1c1a79b7f7f8600fcae8a743dbc60153cc`)
- google.golang.org/grpc@v1.84.0/internal/resolver/dns/dns_resolver.go (source SHA-256: `fae3828426c9cb90d31a0feafce53892434cfce8ec4ceb2f5cd7f823c2d67136`)
- google.golang.org/grpc@v1.84.0/internal/syscall/syscall_linux.go (source SHA-256: `91096efdaa61581ac7b9d2d8076635d8ca5ac7ef3064094dab2f8efb70f003b2`)
- google.golang.org/grpc@v1.84.0/internal/syscall/syscall_nonlinux.go (source SHA-256: `3b74e0c3889ed8a826963f5f3bf9e7bf2bab3ee05f6d2739c9f1bb75cca8d490`)
- google.golang.org/grpc@v1.84.0/internal/transport/defaults.go (source SHA-256: `ca8b8cee5cb10ea5e3396a4e058f31d3a7d7f8e46ce95a82d84895fc34465888`)
- google.golang.org/grpc@v1.84.0/resolver/dns/dns_resolver.go (source SHA-256: `2c2a77109de70df76407a8d6f4dd5153409bca904067fa021e65b31303facd2a`)
- google.golang.org/grpc@v1.84.0/version.go (source SHA-256: `3aea72eb27010cde5430515214c5ee820354a190f725b260d4086f861677bfa4`)

````text
/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/attributes/attributes.go (source SHA-256: `7607fcd184cfdb2a09c1c461477ca83953705a397830287df642d7655f12ea6b`)
- google.golang.org/grpc@v1.84.0/backoff/backoff.go (source SHA-256: `79c43f7f8f0ea8a1aa8a93e0f37a2d9eedd3a6dafde0adee0a0a97a6a0ed5d8d`)
- google.golang.org/grpc@v1.84.0/preloader.go (source SHA-256: `7c05bfee179f64513d49823bf2ede803dfdad8e8f20a819467771c058d401f73`)
- google.golang.org/grpc@v1.84.0/serviceconfig/serviceconfig.go (source SHA-256: `7cd30c7184690c8ed10a42314f53bb857f8a731df83a53ded43848ea73a913dd`)

````text
/*
 *
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/balancer/grpclb/state/state.go (source SHA-256: `ad2437a3ac57666dd96ea8d76d32c52f7ef66bdee59efd0faa2ab2b780805d05`)
- google.golang.org/grpc@v1.84.0/channelz/channelz.go (source SHA-256: `15cd2cc89bf22b6990e5621febbf6f01387b9f6fb29790dd59a1aa60c05eba6a`)
- google.golang.org/grpc@v1.84.0/credentials/insecure/insecure.go (source SHA-256: `956c4bc84a6809e5540f13bb687b0450d3424bc19f695734be4078300be30b6a`)
- google.golang.org/grpc@v1.84.0/grpclog/component.go (source SHA-256: `2b2a50fdc442ef1a92d5012c9f21763b71d47794114604b4dfcb39c410d9c4bf`)
- google.golang.org/grpc@v1.84.0/internal/channelz/logging.go (source SHA-256: `60832e30dae6a10abdd6ed7266f458c20c8b065bd8a211d59f5c8139fb47380c`)
- google.golang.org/grpc@v1.84.0/internal/credentials/spiffe.go (source SHA-256: `149056461f3e6a44cb6f24c7f7f7e836a23103f8ab68336f0b8d0e868da5f907`)
- google.golang.org/grpc@v1.84.0/internal/credentials/util.go (source SHA-256: `4d92ceb2d50747ae531274f5cbde37510711ecaaa4e35c6f88903263e90d6de4`)
- google.golang.org/grpc@v1.84.0/internal/envconfig/xds.go (source SHA-256: `06e17222e05d474eb6a2564a4aa9e2f8ee1169e49ab85252b1bd7a1f4afdffa8`)
- google.golang.org/grpc@v1.84.0/internal/grpclog/prefix_logger.go (source SHA-256: `2a2dd5790da1919601dbf70d49daae0afa576139f5d0d7b6e2f42e29306b70fb`)
- google.golang.org/grpc@v1.84.0/internal/grpcutil/encode_duration.go (source SHA-256: `484fefb3ce4e788e825c6a8e75346633a4059eb45d56a78f79f641735f5b1653`)
- google.golang.org/grpc@v1.84.0/internal/grpcutil/metadata.go (source SHA-256: `a141ed371944b5f4edd385c3b3138283543d166e5b6cf90d1b33987f702a885e`)
- google.golang.org/grpc@v1.84.0/internal/metadata/metadata.go (source SHA-256: `dc817c99dcf0a4755b6a7f945916c775d103e8711a8b70be65cfdcd158afd53a`)
- google.golang.org/grpc@v1.84.0/internal/resolver/config_selector.go (source SHA-256: `fd662a50b21d7d76bfeb0c59a05e2e3247788b5524ef3285b5dbe4b1940cd78c`)
- google.golang.org/grpc@v1.84.0/internal/resolver/unix/unix.go (source SHA-256: `d7413035609466572b9e656055553df55d2f2cdda7c5c35c419877615a51dc38`)
- google.golang.org/grpc@v1.84.0/internal/serviceconfig/serviceconfig.go (source SHA-256: `ee1eb07030513d275455f66d92261a692cb0fd0ddb32371e98a18e72d291caa2`)
- google.golang.org/grpc@v1.84.0/internal/status/status.go (source SHA-256: `687cb347d0538c44b87ce0729b69afef0964692b20a4823a9811363c0d26cfee`)
- google.golang.org/grpc@v1.84.0/internal/transport/networktype/networktype.go (source SHA-256: `f3ad9189353878367ea40a18904b7bca73352dac63180f9d04051cfc16974e84`)

````text
/*
 *
 * Copyright 2020 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/grpcutil/grpcutil.go (source SHA-256: `eaedb035765c70f150d078df3436856ecf9109d5bb72807462a7ce0d3f8cbd49`)
- google.golang.org/grpc@v1.84.0/internal/pretty/pretty.go (source SHA-256: `64e5da0da1abeefe8b1f75dedbc6054b2e98b7dff430b7942f7807e76f6b0486`)
- google.golang.org/grpc@v1.84.0/resolver/map.go (source SHA-256: `38a8a56bdd55fb00474d40d042624f8dbe58697fe2846d26851007a7a8ade903`)

````text
/*
 *
 * Copyright 2021 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/balancer/conn_state_evaluator.go (source SHA-256: `a380fe0f58e11f4a13184768123da0a3eb4cf41a50b838b6eedab0ce5b2f9156`)
- google.golang.org/grpc@v1.84.0/internal/balancer/gracefulswitch/gracefulswitch.go (source SHA-256: `f44b3da3c81ef06997fb2736c1d768fdc69234f83116b5b1359430fb75fba6b0`)
- google.golang.org/grpc@v1.84.0/internal/envconfig/observability.go (source SHA-256: `57c20c861a81a343a7e0f65d29952f5a6429dc0246fbd76fe3f69c9d214ce000`)
- google.golang.org/grpc@v1.84.0/internal/grpcsync/callback_serializer.go (source SHA-256: `c24911cca29f8e47d4c12a08d14bc257d1ae1e57e86b899790a9d4ac457fbc57`)
- google.golang.org/grpc@v1.84.0/internal/grpcutil/compressor.go (source SHA-256: `c574024c853b41afc612b597ad168ed6a680d09fa60009d93842b1897ec9f88c`)

````text
/*
 *
 * Copyright 2022 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/grpcsync/pubsub.go (source SHA-256: `1195a2aa024d126559c8feeddb25f5a79e950560fae5cdb25bce99f2a7184a3d`)
- google.golang.org/grpc@v1.84.0/internal/idle/idle.go (source SHA-256: `f67f127edd0ce323ba14bce05e5f6c4f61bf66c1f1e5c075fd4576898cf491de`)
- google.golang.org/grpc@v1.84.0/internal/resolver/dns/internal/internal.go (source SHA-256: `360c14f956b0fc34a4ccc067945367ad0c16b192dc050f99ab2a6f1f4e7e5c95`)
- google.golang.org/grpc@v1.84.0/internal/serviceconfig/duration.go (source SHA-256: `285f9de62fca3d050acf942df8f207ac5233e7bd6e2aca11ee780ec3728418d1`)
- google.golang.org/grpc@v1.84.0/internal/transport/logging.go (source SHA-256: `ebf772b0c69a9a96b5bc58670c3b4197d707fb4980defca772ac51e492ee6a5e`)

````text
/*
 *
 * Copyright 2023 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/balancer/endpointsharding/endpointsharding.go (source SHA-256: `0c86aace9630a05af24a171d88ea17008f9a829f3b8e6c3b5379904173e4f53f`)
- google.golang.org/grpc@v1.84.0/balancer/subconn.go (source SHA-256: `f3308d7226e5061977fbb1251a2093d77816a2220bae15a531e1406e64b7ee59`)
- google.golang.org/grpc@v1.84.0/encoding/encoding_v2.go (source SHA-256: `216666c3ab2394d582154a5e07ca415fa424dbb6a01cc9e309c188c5b50403f9`)
- google.golang.org/grpc@v1.84.0/encoding/proto/proto.go (source SHA-256: `e49a349919237058e91aa19cf2128f78fdbdaae0f8ed169a600ffcc07ff566bd`)
- google.golang.org/grpc@v1.84.0/experimental/stats/metricregistry.go (source SHA-256: `bd1c67f63b659225510f16afa534f14dea5a272c615a38e0085f4d04ea3d7ef4`)
- google.golang.org/grpc@v1.84.0/experimental/stats/metrics.go (source SHA-256: `e451c55b40035e873f4e8b93e2ae6c2fe9f29c032ff7e2bea2c97d5e1b3a4661`)
- google.golang.org/grpc@v1.84.0/grpclog/internal/grpclog.go (source SHA-256: `ca8b6f4b1662bcc98929293f22f4d3fdab23b32975f7387b6254f5a70df0b985`)
- google.golang.org/grpc@v1.84.0/grpclog/internal/logger.go (source SHA-256: `047e4dc02b8b3b1fdba6f98ce4664f84e43d6476b83b5356f149c6237a3e7eb4`)
- google.golang.org/grpc@v1.84.0/grpclog/internal/loggerv2.go (source SHA-256: `73624458b062c9a9b60b89b0640d5958c53d2703d032dfedbcfed5b68dcd6b37`)
- google.golang.org/grpc@v1.84.0/internal/balancer/gracefulswitch/config.go (source SHA-256: `efa21a7a940f89cdbbb011b6de3b768ac6d8a51967dfc6b5115c7baa755985cf`)
- google.golang.org/grpc@v1.84.0/internal/channelz/channel.go (source SHA-256: `56d6dced026aa9dcd96a862cb62c29f6a24ded63afaa5192ec148af3de7c84ca`)
- google.golang.org/grpc@v1.84.0/internal/channelz/server.go (source SHA-256: `02e96d7746fef4d0dccaa601f021a51554e6ab8b1e39b8584c295a65db2a0dbb`)
- google.golang.org/grpc@v1.84.0/internal/channelz/socket.go (source SHA-256: `ac4c79498ec3fa4364c0dd55a5ccd263bef4374cc7f2600dc01f3e77a0401044`)
- google.golang.org/grpc@v1.84.0/internal/channelz/subchannel.go (source SHA-256: `290dc50897fb7021ac877753614cb3e6051241860b89ae286cc3142facf41863`)
- google.golang.org/grpc@v1.84.0/internal/proxyattributes/proxyattributes.go (source SHA-256: `567a57b2d6a5795a028cf8ee97809264a8cbf1680690d53897a95ea6225d4d69`)
- google.golang.org/grpc@v1.84.0/internal/resolver/delegatingresolver/delegatingresolver.go (source SHA-256: `3c23d9fee64ae948fc47613489dce867eeb60c5bae8e6e5ba532db56b92ef545`)
- google.golang.org/grpc@v1.84.0/internal/stats/labels.go (source SHA-256: `a768974a8542c4df9987bbb582e9fc79e032b01a544f170942e6d8f3f94a38a2`)
- google.golang.org/grpc@v1.84.0/internal/transport/client_stream.go (source SHA-256: `813aafffcf7cf6310d145379218c020b256558947b4f4741a762417657106874`)
- google.golang.org/grpc@v1.84.0/internal/transport/server_stream.go (source SHA-256: `aa8d18fb85ba760bef634e337cca4160e4b221808bb9f777fb87841b1457eda2`)
- google.golang.org/grpc@v1.84.0/mem/buffer_pool.go (source SHA-256: `1d5943e9b040303eb71146b513a4fb2c9630bf047a3d064c70aa6dafd735846a`)
- google.golang.org/grpc@v1.84.0/mem/buffer_slice.go (source SHA-256: `eb963947548f5a64d07df1f79d06746adfdb597cef66e6b4207ff4389edd8340`)
- google.golang.org/grpc@v1.84.0/mem/buffers.go (source SHA-256: `77fdd82c81c40258eb67e20cf83f42c1ae6f30191a5f27634b59c300a27ebb4d`)
- google.golang.org/grpc@v1.84.0/stream_interfaces.go (source SHA-256: `c51281c4729c8a806e038104b218ef8cf2eb2c90012c35cc6e0f107e603d6e38`)
- google.golang.org/grpc@v1.84.0/trace_withtrace.go (source SHA-256: `874bdd0bb0a18716d92d8bc131e2ed5e0fce08fef3025a3e1f8f5c91a6842dfc`)

````text
/*
 *
 * Copyright 2024 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/encoding/internal/internal.go (source SHA-256: `2d4c22473d74a56aa0293ffc042c01c1f4a233e1ea375c851b82ca23f2755c6d`)
- google.golang.org/grpc@v1.84.0/experimental/balancer/weight/weight.go (source SHA-256: `aa4a1308735237163112c0c450e2cc5eba2bc18f847f31a0a586d29db29466c0`)
- google.golang.org/grpc@v1.84.0/internal/stats/stats.go (source SHA-256: `d026ac30aa6c235e50dd6c9e76c600e177b7c7adfd7a2ed0845d375f4f0bdd26`)

````text
/*
 *
 * Copyright 2025 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/clientconn_disconnect_reason_noplan9.go (source SHA-256: `23f1d780ca5185a832b94dc81111af6ea8916360631b450268adf288fad0e68a`)
- google.golang.org/grpc@v1.84.0/internal/grpcsync/refcounted.go (source SHA-256: `381c21e880738215402668bddd6a692e439c48cb27c25554e5641692056a1c29`)
- google.golang.org/grpc@v1.84.0/internal/mem/buffer_pool.go (source SHA-256: `ffcb191777909f96dd1677977c9fef2f897cca5243439f76780d7ddc576ad92a`)
- google.golang.org/grpc@v1.84.0/internal/transport/internal/internal.go (source SHA-256: `5429a838d99d91836b457d3832b6e501ef3dc318310a15438953be65593c3c1d`)
- google.golang.org/grpc@v1.84.0/internal/transport/readyreader/raw_conn_linux.go (source SHA-256: `05a068012fd512ce4ad3c5471f9546349b4e9acdad30dde5df0f02917af58d3f`)
- google.golang.org/grpc@v1.84.0/internal/transport/readyreader/raw_conn_nonlinux.go (source SHA-256: `b9f0cf777d60fa3acfb2c22d263738a361b30ed484885b039db1c555a8f53256`)
- google.golang.org/grpc@v1.84.0/internal/transport/readyreader/ready_reader.go (source SHA-256: `aefa4f5b55e10c0ff86d5cf5b5d77a44eb79ca00c412dc9218087f23537fdb00`)

````text
/*
 *
 * Copyright 2026 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/internal.go (source SHA-256: `ac4bdac04665525268857f3ac54386f1d9a0aeb6c2f41b11be4d4a7a2799669d`)

````text
/*
 * Copyright 2016 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/buffer/unbounded.go (source SHA-256: `58a4721d17960274119fea6257cb09dc362fd186b23a6f96711686961674d6eb`)

````text
/*
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/balancerload/load.go (source SHA-256: `6678fb7a22285da1d63ab9d8a7a98f795be776506af243224748253bd137ab61`)

````text
/*
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/credentials/credentials.go (source SHA-256: `1def08c84b6dac9fe9030d3a78d46293ca9549ac40b3ec8e91308369aba5dcdf`)

````text
/*
 * Copyright 2021 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/experimental.go (source SHA-256: `cee5c034eb7a26d6d20fb30105ac8961b450ee4082dd9b9e7a3cd37b83029610`)
- google.golang.org/grpc@v1.84.0/internal/tcp_keepalive_unix.go (source SHA-256: `e8bfe03234b391d24006a3a274590111f0f8705fc5b25d9a78391bfdde3df32c`)
- google.golang.org/grpc@v1.84.0/internal/tcp_keepalive_windows.go (source SHA-256: `55f9a0b5c7d247ed8415f0e907877e4e67078ae9e9b571a999ad7301c8cd45f2`)

````text
/*
 * Copyright 2023 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/balancer/pickfirst/internal/internal.go (source SHA-256: `2abb7cfbfe55770e52f54381da1e426a173eccb663cc4c9b070549087c7a7ca9`)

````text
/*
 * Copyright 2024 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/internal/stats/metrics_recorder_list.go (source SHA-256: `e21ee4e6827c8cda9d20677f1eacb7721ee253785f87b475b58b3af609f62d40`)
- google.golang.org/grpc@v1.84.0/stats/metrics.go (source SHA-256: `4dd20d08cc70a7d586d31b866c1e82df930764ed593efeec149ecde448b0e747`)

````text
/*
 * Copyright 2024 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
````

### Original source notice

- google.golang.org/grpc@v1.84.0/binarylog/grpc_binarylog_v1/binarylog.pb.go (source SHA-256: `79943c635d176a7cfacb6f4e327df3c754378622a9eee03edea41d4ae4adc211`)

````text
// Copyright 2018 The gRPC Authors
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
````


## google.golang.org/protobuf@v1.36.12

Source: https://pkg.go.dev/google.golang.org/protobuf@v1.36.12

Targets: darwin/amd64, darwin/arm64, linux/amd64, windows/amd64

Module checksum: `h1:pJOKDDOyeXErUroCihFAd5LQuwXBSpVnKGrj5o/fwxc=`

### google.golang.org/protobuf@v1.36.12/LICENSE

SHA-256: `4835612df0098ca95f8e7d9e3bffcb02358d435dbb38057c844c99d7f725eb20`

````text
Copyright (c) 2018 The Go Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

### google.golang.org/protobuf@v1.36.12/PATENTS

SHA-256: `96f408bfae65bf137fc2525d3ecb030271c50c1e90799f87abf8846d8dd505cc`

````text
Additional IP Rights Grant (Patents)

"This implementation" means the copyrightable works distributed by
Google as part of the Go project.

Google hereby grants to You a perpetual, worldwide, non-exclusive,
no-charge, royalty-free, irrevocable (except as stated in this section)
patent license to make, have made, use, offer to sell, sell, import,
transfer and otherwise run, modify and propagate the contents of this
implementation of Go, where such license applies only to those patent
claims, both currently owned or controlled by Google and acquired in
the future, licensable by Google that are necessarily infringed by this
implementation of Go.  This grant does not include claims that would be
infringed only as a consequence of further modification of this
implementation.  If you or your agent or exclusive licensee institute or
order or agree to the institution of patent litigation against any
entity (including a cross-claim or counterclaim in a lawsuit) alleging
that this implementation of Go or any code incorporated within this
implementation of Go constitutes direct or contributory patent
infringement, or inducement of patent infringement, then any patent
rights granted to you under this License for this implementation of Go
shall terminate as of the date such litigation is filed.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/encoding/prototext/decode.go (source SHA-256: `6838aa95ca535e4bb7c331459b899ce75a75b0299a64579b6ec8524dbe3515f1`)
- google.golang.org/protobuf@v1.36.12/encoding/prototext/encode.go (source SHA-256: `a139182ec75da9342d6871316af06633c24d0f75317b92ebdf792f41753a6a91`)
- google.golang.org/protobuf@v1.36.12/encoding/protowire/wire.go (source SHA-256: `109e891e95fb059a3622e2ad97bb06f6829da7d4ac7fc13c36deb508d6f9ef40`)
- google.golang.org/protobuf@v1.36.12/internal/descfmt/stringer.go (source SHA-256: `d54f3a4bb048f1e62b4d49274f313417f48451c796c309542aa9b26180a1be90`)
- google.golang.org/protobuf@v1.36.12/internal/detrand/rand.go (source SHA-256: `c01386e38186772e042196f4b03ce25eed8993ea8ffb8c012f43a2286b7752a5`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/defval/default.go (source SHA-256: `2ec7a1ef4aa6710853d0f7003310f6808ea804adc0c495021814ca7541cdbb80`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/json/decode.go (source SHA-256: `149097bfd4ae2cd52f8f8d5b36c850c3972124670c75aaabfaa32edb29512aa4`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/json/decode_number.go (source SHA-256: `ae8d84879231dbbe42065f7d3852357046382ea3e290c8262fb475559638f2e0`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/json/decode_string.go (source SHA-256: `99540744259eb12d78fc7431c506dcba4c4148eed66c21037b7a860a940bb6cd`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/json/encode.go (source SHA-256: `ce3d1d9568fa93f2e0b383453198bf6e0c75ef4ccf0e84141622413d6ccb693a`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/tag/tag.go (source SHA-256: `5cbf357c44751276c273ee8e4d46851acc7a55579b742b9cbadaf29e142077de`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/text/decode.go (source SHA-256: `59d2e1a887cbe3fd657c68a65984fd3255cfb43de1dde6e529161f6959ac74f6`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/text/decode_number.go (source SHA-256: `1bb6f72d94ed0866666810bb23664286ba28916092e11df9679bdb0c32ad98cd`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/text/decode_string.go (source SHA-256: `5974107818845d7d94c39cd275985dfbcd55b9eba4d731e182715f985370e7ba`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/text/decode_token.go (source SHA-256: `b51ba5d349f6c36f99544020733af26dbcb476bd73548cd943c14286fae2c475`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/text/doc.go (source SHA-256: `1ca496689a8d7d06862f6b9dab7982d1e751da4c6e97ae59b747890b7a5c541d`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/text/encode.go (source SHA-256: `1a906e4182f13d2120f3af53aa2277cc131449937f545b5bd5e7e340c7bcf606`)
- google.golang.org/protobuf@v1.36.12/internal/errors/errors.go (source SHA-256: `da78e47d269681fea7a751c5f4e1874abc9cce354011be7f5f06a4d4656f3632`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/desc_list_gen.go (source SHA-256: `5157ced1a76ed9eb6dcf4360263b3959d383c2b62c14a0fc2d06a1979d0067d9`)
- google.golang.org/protobuf@v1.36.12/internal/flags/flags.go (source SHA-256: `36bec6da6332d7dc34aeb7b81a5b77c5b1aa71be761bbfd5d341f519a9de921c`)
- google.golang.org/protobuf@v1.36.12/internal/flags/proto_legacy_disable.go (source SHA-256: `c083481c9b7266f6b7c78dd1be24147428f6e5677b58c479d742586bd9f2eebc`)
- google.golang.org/protobuf@v1.36.12/internal/impl/api_export.go (source SHA-256: `688c3577d3439d9b5aef2e5c42a875b8074117b7c2f506b4f74a6c3484ddb218`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_gen.go (source SHA-256: `4eed98966d930db1855cadf0972d255cd89733e29815ef4d3904e3fe0e473a25`)
- google.golang.org/protobuf@v1.36.12/internal/impl/convert.go (source SHA-256: `2e5bca1cb7e7dee3b21f2d9344b54c748f75a17d2ce9214ded04e05db0f3e81b`)
- google.golang.org/protobuf@v1.36.12/internal/impl/convert_list.go (source SHA-256: `2eb7936a7cad5faaa93c0da719652c06cd2edd1a9c5876641ff3f2f580588c86`)
- google.golang.org/protobuf@v1.36.12/internal/impl/convert_map.go (source SHA-256: `dc384760b90c4300f8a3e5377840b7dd42bd3596c9849887021a1bf75a000b73`)
- google.golang.org/protobuf@v1.36.12/internal/impl/legacy_enum.go (source SHA-256: `32b4b7fb059c4f8189cb3e6e8d4b5a68f35380a3efae93499f595f5074fa6471`)
- google.golang.org/protobuf@v1.36.12/internal/impl/legacy_extension.go (source SHA-256: `2b07af9c4eaf0b7ad439b6f1937f65d2d33bb9f1691486e6d4d73b45331a2118`)
- google.golang.org/protobuf@v1.36.12/internal/impl/legacy_file.go (source SHA-256: `d94f1e0dd461f2be37fc6c44bf40f794114ed6d52112afe640a298541c66e9da`)
- google.golang.org/protobuf@v1.36.12/internal/impl/legacy_message.go (source SHA-256: `a7f558df997f1e7b5c29210281de1f9a6b03c1afbf6b0ce18fc27dc366b5c9ac`)
- google.golang.org/protobuf@v1.36.12/internal/impl/merge_gen.go (source SHA-256: `8eddd46272d1cf5e7de3e55f1c92646ad0fc1e2be9d5929633aa5b0813b5780e`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message.go (source SHA-256: `6613d5b02b1c4f29f719d4bf6765e6c70015b63a1536679a2939b49823caa1ac`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message_opaque_gen.go (source SHA-256: `ccb4dea1ba262a75ac315fe29781579d8cc69fc94c48b9f2b955ac5e3eb50309`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message_reflect_field.go (source SHA-256: `249b98128b4a88b48c7693465be56af448f5ff62e2a0b1def1ded23cf06e1e37`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message_reflect_field_gen.go (source SHA-256: `b502db57feef273931b0508b8d863ee4ca6a1576ac2a2d66cf740d14a1a56e75`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message_reflect_gen.go (source SHA-256: `4478a28ca4a028b25618ed03c5cca1a2c15dd7f55c64799683ca661de8d0d901`)
- google.golang.org/protobuf@v1.36.12/internal/impl/pointer_unsafe.go (source SHA-256: `e83e0d25bcc8eb9a11ea27e2f8f433cc5b2e5ff5e44d3d54bc7940e50f966133`)
- google.golang.org/protobuf@v1.36.12/internal/pragma/pragma.go (source SHA-256: `d3babf1e4b0b32f721a0dd3e6c3e31fc573d550fd9a1d51e543f40482b0fc6f5`)
- google.golang.org/protobuf@v1.36.12/internal/set/ints.go (source SHA-256: `5876c26b3c816d4f662b737731e89d3a44317c798c6c5400526c32387292e16d`)
- google.golang.org/protobuf@v1.36.12/internal/strs/strings_unsafe.go (source SHA-256: `35b1b1589491aa6c09707f4442bafded55e95c4972a14c64affa248b7bb57709`)
- google.golang.org/protobuf@v1.36.12/proto/decode.go (source SHA-256: `5259b0da78b384fcce573d774277f13e83a1f06b7fd212dd5f6fc617b3941917`)
- google.golang.org/protobuf@v1.36.12/proto/decode_gen.go (source SHA-256: `d8581139cb95720ebe63a1a5654dd1d5927543c2b1ad1ef37ed37ac41e7f8056`)
- google.golang.org/protobuf@v1.36.12/proto/encode_gen.go (source SHA-256: `8a5d8bc917c21a19c1c5ef4b58f747880cfe1321de1294662d0ee2c903c7dc9e`)
- google.golang.org/protobuf@v1.36.12/proto/proto.go (source SHA-256: `8a87120799a6f0ed42a9096cd997d93baeaf143db660daa5b9b8c6e0665d95ec`)
- google.golang.org/protobuf@v1.36.12/proto/size_gen.go (source SHA-256: `e8def0182318c353c022c973eec36615c121078ac405b059c7f02fad71c86b8d`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/proto.go (source SHA-256: `87783d5e1fedf825ec7ff63c36494e97011e221c466fc81affaacac618fdb907`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/type.go (source SHA-256: `1f0cab06a41ac3b25fa02512b0dbe94acb62ad9b121bd8739d9f4615d3ae8691`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/value.go (source SHA-256: `78cb6818d7abba4028bc152bf273c9369c50a0f82c5faa248e1f7a5fa31da385`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/value_union.go (source SHA-256: `efa2f612fa2bbaac079d7c6b1718fa50ce915488c6d8fd2adde883bb3d47b7ed`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/value_unsafe.go (source SHA-256: `fd19f2ed3fc1c3c70a1d2074c2867e2ae47947c583645172496c0257f1e72c2f`)
- google.golang.org/protobuf@v1.36.12/reflect/protoregistry/registry.go (source SHA-256: `75699e3c2cc3650473c04819954673e7382a101f6f7cbc11a30f8bec6fabba0e`)
- google.golang.org/protobuf@v1.36.12/runtime/protoiface/legacy.go (source SHA-256: `9d79e425272f9fa102c269f4cb7f8cbe405d0f2a9af075e150235df30b3bd617`)
- google.golang.org/protobuf@v1.36.12/runtime/protoimpl/impl.go (source SHA-256: `b47b677332e32d41f1c3a71dabfe8c73945d62823f81d39f1045d1091f38f858`)

````text
// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/encoding/protojson/decode.go (source SHA-256: `240f38a45c0b0b9000be8e0862b1d8d720cf6ed245ba77cb95cbe641fa4dfc2c`)
- google.golang.org/protobuf@v1.36.12/encoding/protojson/doc.go (source SHA-256: `b286de630dedcb17f1d2369f7faf598fa73ab40aa284ea27171cdb164889e847`)
- google.golang.org/protobuf@v1.36.12/encoding/protojson/encode.go (source SHA-256: `76d996dd40a3cf689b49d8c787ee857388277e7a3bd45b4bceec93dc8fabdee1`)
- google.golang.org/protobuf@v1.36.12/encoding/protojson/well_known_types.go (source SHA-256: `811976bfadc6497cd00b6cb5d1fd6ec6a94d3f39118b177f7b6ee0894f497636`)
- google.golang.org/protobuf@v1.36.12/encoding/prototext/doc.go (source SHA-256: `48657fec4cff435f30956a1dc6c1b639550ed5fffe1633f1c976645c67d199e7`)
- google.golang.org/protobuf@v1.36.12/internal/descopts/options.go (source SHA-256: `0787d3db599ef201cc7e16598912030580be7023cde1da4bdf52c60eda3d2309`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/json/decode_token.go (source SHA-256: `1b60ca4b6380800bf748d2dfbf8c7b0dff8f52c8f223d3fe02c4879b2072f6b7`)
- google.golang.org/protobuf@v1.36.12/internal/encoding/messageset/messageset.go (source SHA-256: `4a54d8cd4c4c0236d9ceb8200580c9e5d13e5157162d0e2654b8de83707b64bb`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/build.go (source SHA-256: `52bef727d5e493d495dc5ceca4eb21801e006bcf46db23367296b64058e12037`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/desc.go (source SHA-256: `418bf5762126ff47c5e03ccde0e6062cef63c0148b56f7197f4a19ddb5bc64ca`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/desc_init.go (source SHA-256: `7a1d253f0b69cf6ae0c2c86d88dd943202228a47e16dd850a9c8b1e7ad66c9a5`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/desc_lazy.go (source SHA-256: `8c326a5637883aa7a44f96be79b34c6000780bf55449b52cbf74b421d57e0f2b`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/desc_list.go (source SHA-256: `78598da1e3c79f00189bba62c274d0d6febd4a2c972b269bbc0451c1218d5310`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/placeholder.go (source SHA-256: `f3fa4a70a986c68e07b5a7e6571428f481f0d8c3832e2fa1ccc58064709dcc98`)
- google.golang.org/protobuf@v1.36.12/internal/filetype/build.go (source SHA-256: `e22ba313e842bd9e8274b5f9dfdec5896bc9712830f1df75809bb8a0f3d06a34`)
- google.golang.org/protobuf@v1.36.12/internal/genid/any_gen.go (source SHA-256: `e52abb175ed291da379fe73582a640b13d08ae426b51f6e40a82fc07ef4494f2`)
- google.golang.org/protobuf@v1.36.12/internal/genid/api_gen.go (source SHA-256: `9b02e61ea3e584147063ee7fd3ffecd946e96e2238da6c838e5262912cce3dfa`)
- google.golang.org/protobuf@v1.36.12/internal/genid/descriptor_gen.go (source SHA-256: `caa818838cb2eba6e9278dabdbe2a1af31f24155904b6b2e5634abbead928ace`)
- google.golang.org/protobuf@v1.36.12/internal/genid/doc.go (source SHA-256: `25acd9d2f70b5384f769932a362d5f28d20a71c86720978401b43452e6fdfc0e`)
- google.golang.org/protobuf@v1.36.12/internal/genid/duration_gen.go (source SHA-256: `223477c7c21d3aba70b6afaad202499d35dd28d97bb3096649d398ec09ffbd20`)
- google.golang.org/protobuf@v1.36.12/internal/genid/empty_gen.go (source SHA-256: `1bb2b55b36ca99633efc148681988dbd738bfc6d09b88c92bed6392a4a16e80c`)
- google.golang.org/protobuf@v1.36.12/internal/genid/field_mask_gen.go (source SHA-256: `97ac7e90145a527f1b53e26868cb2257af5ca000103dcf20ff6374693512d5de`)
- google.golang.org/protobuf@v1.36.12/internal/genid/go_features_gen.go (source SHA-256: `e13a8023f516a885854ebae3f417ecf2bc55da9bbb2385b38fcf8b91010b6cc6`)
- google.golang.org/protobuf@v1.36.12/internal/genid/goname.go (source SHA-256: `78d5d5dc16c62d8990f1dfaef2d744aff13cdae8aa3f03595813b374ae9743d5`)
- google.golang.org/protobuf@v1.36.12/internal/genid/map_entry.go (source SHA-256: `539c26c0b3ab5e42ccc3235ea14f3b82ccf1be638e2e1a6602cfc2e133a57201`)
- google.golang.org/protobuf@v1.36.12/internal/genid/source_context_gen.go (source SHA-256: `e0aac89c39dea4892ba50ea5316728a8c582ba9a1434406c5aa047c0e90bcdd0`)
- google.golang.org/protobuf@v1.36.12/internal/genid/struct_gen.go (source SHA-256: `df9f6562b19e5c60ac1e6006365a7ddd0c7fcca2101dd1da909413c0cd0ef6dc`)
- google.golang.org/protobuf@v1.36.12/internal/genid/timestamp_gen.go (source SHA-256: `5652b719e792955d360f4b5c8f8889923976b55f86e13dd233edf366fdd71e2a`)
- google.golang.org/protobuf@v1.36.12/internal/genid/type_gen.go (source SHA-256: `13bda96c10ef3bdbfcd81a1e895fdf3bb9d071714eade6a86ef4a1d09b1505bf`)
- google.golang.org/protobuf@v1.36.12/internal/genid/wrappers.go (source SHA-256: `a68cd5fd631dcac3ac053c78cdecefd24f909aae9920da54e801e0a496af89dd`)
- google.golang.org/protobuf@v1.36.12/internal/genid/wrappers_gen.go (source SHA-256: `eaa8ad201d3ef370a565fea90ae9f5bef6efec40c7fa58a53e7efe635731e887`)
- google.golang.org/protobuf@v1.36.12/internal/impl/checkinit.go (source SHA-256: `6d1b9b5da1fad3d8cc5efc939e90672dc1690df3fa3f296e140685d78d7167eb`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_extension.go (source SHA-256: `4bbee1d9036b26a17b1b3279c56c59b7d5ac16b65098ed65011ee3aba4de3ae9`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_field.go (source SHA-256: `d794d96c2d98db80f7230cebdd8f5b118d6ecd31e08ac445689b239dd110f65a`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_map.go (source SHA-256: `85f1d364d793344cf724b7204110f79c77d9a75e1b2d3e8ced3a7a0a42f7aa9a`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_message.go (source SHA-256: `0533a067562e6181ed035fb9bff3741f264bd9a671408df06e3d059b8d61d49a`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_messageset.go (source SHA-256: `f0d06fd5031c23940ee5b4ad63a7158d8e883c0c0dde56e0f0049be69dcdb3b7`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_tables.go (source SHA-256: `8d4cdbd1367e0c089ecada4431bdb0b2f3198aa3907c25887131650adf78d29b`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_unsafe.go (source SHA-256: `ac4d881c0b561071e810302fb075c5de804e947f33cf4e0901f8876590e84c29`)
- google.golang.org/protobuf@v1.36.12/internal/impl/decode.go (source SHA-256: `12b24500e3708913ecba846287a84f413e36ca496bfb713b931102e0041621d3`)
- google.golang.org/protobuf@v1.36.12/internal/impl/encode.go (source SHA-256: `bec7e80f5cda21616d92312215d4d0c4b11ce891ccc2e4d16447259dbacc9739`)
- google.golang.org/protobuf@v1.36.12/internal/impl/enum.go (source SHA-256: `0501ee078196bb8916290dcf2eb4d02f1e3f9354f7ef034e4be0695a4e13c9b5`)
- google.golang.org/protobuf@v1.36.12/internal/impl/extension.go (source SHA-256: `4642ab519bca05f2d4512d3a80b45ba3b09f1ccb6547faefb71d0f1d75e121ed`)
- google.golang.org/protobuf@v1.36.12/internal/impl/legacy_export.go (source SHA-256: `e338223a953cb22d87c0621d7e6c373d4a39853c98db4999711ea412b2d94d3f`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message_reflect.go (source SHA-256: `3250ea5ca291396cca8c7496aec863db78033f2a57fef8f4b2d2f357f5f2180b`)
- google.golang.org/protobuf@v1.36.12/internal/impl/validate.go (source SHA-256: `70628f96f7c80a36f944234dad217660d74efe59f1cf299f4430a2145d5ac548`)
- google.golang.org/protobuf@v1.36.12/internal/strs/strings.go (source SHA-256: `97aaee23976801265248c5462662e5061a2af09956c4d6b12990e16b5acf0c8d`)
- google.golang.org/protobuf@v1.36.12/internal/version/version.go (source SHA-256: `1f7c8c097f20bc3f9eb689c37f63151b813ff1f7bf60db296f3a2bf1d77fcaf7`)
- google.golang.org/protobuf@v1.36.12/proto/checkinit.go (source SHA-256: `f95547595ee464a0d2f2bfae729afb61c91c5ed3bc64f089949b4b1555a71b9b`)
- google.golang.org/protobuf@v1.36.12/proto/doc.go (source SHA-256: `3df50c4bae0b7ba05913716e06bff9b1c90f17b2543c6e20a20702d2d62add5d`)
- google.golang.org/protobuf@v1.36.12/proto/encode.go (source SHA-256: `e9a13c44a81991067d003f235234fa8a859ab1f78ff0430212861f5027d5a434`)
- google.golang.org/protobuf@v1.36.12/proto/equal.go (source SHA-256: `d0aa54ea8ad15305c5c8117f86e7c79fbc3e87eb7c62290823b920bb8f902e88`)
- google.golang.org/protobuf@v1.36.12/proto/extension.go (source SHA-256: `4e8b480347db15bee7878d629f22b20782b9a7a3b5f627ffe1fc5ec9f927234a`)
- google.golang.org/protobuf@v1.36.12/proto/merge.go (source SHA-256: `f5c50384cb64947ee71b6e77633c4c88f277575b36121bd7c8ee9f38fb6a16b7`)
- google.golang.org/protobuf@v1.36.12/proto/messageset.go (source SHA-256: `539fb7c4cc85c07a8265ba7b1c85f28ae0292f74ee3ffbdb4a259eb4d1b37273`)
- google.golang.org/protobuf@v1.36.12/proto/proto_methods.go (source SHA-256: `4a5be7ac829abca1043c568c33af31cfac723c2fe688dbff34a605954e5ce141`)
- google.golang.org/protobuf@v1.36.12/proto/reset.go (source SHA-256: `45097aa639f3c64a5533572cb1968d3870c985b3e0b021ce5c0c3b4934768604`)
- google.golang.org/protobuf@v1.36.12/proto/size.go (source SHA-256: `1be720967c0c319c4eb4b438ae37968c74863c6da42caf45ca5b241261ba1c0e`)
- google.golang.org/protobuf@v1.36.12/proto/wrappers.go (source SHA-256: `7068ed0c35932f990aa75bf00d30f6a2d2a3c2c89ebd14ac07aec523013fafb4`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/source.go (source SHA-256: `645962b399787a757fea01ad489c353d11cbe796df1231a89b9dd6fe53c79830`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/source_gen.go (source SHA-256: `a3e4aa29534e9ff9307a73738225c88764e9477271bbf8acac0e6fbf38f4822f`)
- google.golang.org/protobuf@v1.36.12/runtime/protoiface/methods.go (source SHA-256: `44f3bdcd93d290e8fb6b39cbac5c397371dded349632e211e8b1f36a5b53674d`)
- google.golang.org/protobuf@v1.36.12/runtime/protoimpl/version.go (source SHA-256: `d447c5756d628bc1b411a46df732be052b6ad7ee4582aa357ab7776db759df4a`)

````text
// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/internal/impl/merge.go (source SHA-256: `5406a959d280825d0156c8e2e2857bd4e943a1aa0476a44717339ecb8d580628`)
- google.golang.org/protobuf@v1.36.12/internal/order/order.go (source SHA-256: `07529bca11498b579b4d02849560f7fb52ee6cc8cbfcfb82ede169728d8ee4a1`)
- google.golang.org/protobuf@v1.36.12/internal/order/range.go (source SHA-256: `aa723c00dd54eb6f95a35322f93e01e3eb247fffe9af2c1ccb4e846ad5a63b40`)
- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/methods.go (source SHA-256: `d2b015d684ff7405dcf44a766e54c27ffc7c9f440949d7776385a39d5f10aa8c`)

````text
// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/reflect/protoreflect/value_equal.go (source SHA-256: `a5b85845963e9968ecca24ffc5d0b8550bce6967587aaab44a4c797c0845ee55`)

````text
// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/protoadapt/convert.go (source SHA-256: `090a3004e742ab79a549844c8ff3adbe2188ca8318c45357cd02b7a6827a40d6`)

````text
// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/internal/editiondefaults/defaults.go (source SHA-256: `d7314e9ee082d2e56e608ace22c4d3866b9e4e5ed23849996e610dca15e8173f`)
- google.golang.org/protobuf@v1.36.12/internal/filedesc/editions.go (source SHA-256: `8ab88e19e94758aac48d1dca44b028185cc39bd51efbaae8a1add20a895d62d2`)
- google.golang.org/protobuf@v1.36.12/internal/genid/name.go (source SHA-256: `760bff2881021b34c6ed15d0f7d405bec7790bcfc9d5456731fa31053f324a2b`)
- google.golang.org/protobuf@v1.36.12/internal/impl/api_export_opaque.go (source SHA-256: `5e39739bc8ac6d099feb6783d3e3645ceb5ba0aa6926b4b75d826068c4e30dff`)
- google.golang.org/protobuf@v1.36.12/internal/impl/bitmap.go (source SHA-256: `8843d650754650a7cfa627c0b0d5a63465a95afbad87daffd7d5877735ba4da4`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_field_opaque.go (source SHA-256: `30d7554836f43a6cc90cc7651b247221b3252f71885dfc4f51d75dba9f548a23`)
- google.golang.org/protobuf@v1.36.12/internal/impl/codec_message_opaque.go (source SHA-256: `bd5ba6036a439a920d66bc8db90e396a9458dc078ccbbf93a77bc85f224edf6a`)
- google.golang.org/protobuf@v1.36.12/internal/impl/equal.go (source SHA-256: `a8c25efd138002017fcd3c43c7c0d60e7558d45daa6e89a303d9bb18ca7917de`)
- google.golang.org/protobuf@v1.36.12/internal/impl/lazy.go (source SHA-256: `90757f90f660fa9437dde49fe6255b9cb53ac58460db607053b0a550802cd54c`)
- google.golang.org/protobuf@v1.36.12/internal/impl/message_opaque.go (source SHA-256: `45eaa82559904f82f42e4447caae72014bdc338f20a6ff60854f46319c755cee`)
- google.golang.org/protobuf@v1.36.12/internal/impl/pointer_unsafe_opaque.go (source SHA-256: `fef041f0f7b3266405bdffc9fefbfce2be413787e8640c75a941e8056595646a`)
- google.golang.org/protobuf@v1.36.12/internal/impl/presence.go (source SHA-256: `96bc1410475292a56be02d2910921879c8b3a6ecdf147bd85ca60f6650f29bcb`)
- google.golang.org/protobuf@v1.36.12/internal/protolazy/bufferreader.go (source SHA-256: `581604e454f0b0f6bd6e1ba4c31fdb3aa540c72d9259169e4f9e16b1dc4e5f5f`)
- google.golang.org/protobuf@v1.36.12/internal/protolazy/lazy.go (source SHA-256: `fc6ae71902e09b41ae7208fa1e78d0f7be31b8199470d55c5b6078fbe28cc469`)
- google.golang.org/protobuf@v1.36.12/internal/protolazy/pointer_unsafe.go (source SHA-256: `28ffb8eee996ef2bd1c8cf98884531694a45e053793d886549cdc8fa629064c8`)
- google.golang.org/protobuf@v1.36.12/proto/wrapperopaque.go (source SHA-256: `1ca00fab0bcfe001f2082a8bafaf2a8133eeff03cc8d93a6e2f4cb7857a741be`)

````text
// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/internal/filedesc/presence.go (source SHA-256: `b11c16b5e983af652329b2765a49a6e2a9abcd68b3489716e0e86a641ecb5fa7`)

````text
// Copyright 2025 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
````

### Original source notice

- google.golang.org/protobuf@v1.36.12/types/known/anypb/any.pb.go (source SHA-256: `0e45ec8af43084d6dd09d4a300d10d3a51a2790535721585f6531cafae922200`)
- google.golang.org/protobuf@v1.36.12/types/known/durationpb/duration.pb.go (source SHA-256: `c40e56dc8f8643a1a9fcaff45563c11d50f679ea41f37dfb9942c9be31ee0f4b`)
- google.golang.org/protobuf@v1.36.12/types/known/timestamppb/timestamp.pb.go (source SHA-256: `7a3dd04fd975b5e834ebbb4966a824b99707b8549ba1b40b363102967ae1054a`)

````text
// Protocol Buffers - Google's data interchange format
// Copyright 2008 Google Inc.  All rights reserved.
// https://developers.google.com/protocol-buffers/
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
````

