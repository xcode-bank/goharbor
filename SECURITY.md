# Security Release Process
Harbor is a large growing community devoted in creating a private enterprise-grade registry for all your cloud native assets. The community has adopted this security disclosure and response policy to ensure we responsibly handle critical issues.

## Supported Versions
This section describes the maximum version skew supported between various Harbor releases. Harbor versions are expressed as **x.y.z**, where **x** is the major version, **y** is the minor version, and **z** is the patch version, following [Semantic Versioning terminology](https://semver.org/). 

### Support Policy
The Harbor project maintains release branches for the three most recent minor releases. Applicable fixes, including security fixes, may be backported to those three release branches, depending on severity and feasibility. Patch releases are cut from those branches at a regular cadence, or as needed. The Harbor project typically has a minor release approximately every 3 months, maintaining each minor release branch for approximately 9 months.

There is no mandated timeline for major versions and there are currently no criteria for shipping a new major version (i.e. Harbor 2.0.0).

### Minor Release Support Matrix
| Version | Supported          |
| ------- | ------------------ |
| Harbor v1.7.x   | :white_check_mark: |
| Harbor v1.8.x   | :white_check_mark: |
| Harbor v1.9.x   | :white_check_mark: |
 
## Reporting a Vulnerability - Private Disclosure Process
Security is of the highest importance and all security vulnerabilities or suspected security vulnerabilities should be reported to Harbor privately, to minimize attacks against current users of Harbor before they are fixed. Vulnerabilities will be investigated and patched on the next patch (or minor) release as soon as possible. This information could be kept entirely internal to the project.  

If you know of a publicly disclosed security vulnerability for Harbor, please **IMMEDIATELY** contact cncf-harbor-security@lists.cncf.io to inform the Harbor Security Team.
 
**IMPORTANT: Do not file public issues on GitHub for security vulnerabilities**

To report a vulnerability or a security-related issue, please email the private address cncf-harbor-security@lists.cncf.io with the details of the vulnerability. The email will be fielded by the Harbor Security Team, which is made up of Harbor maintainers who have committer and release permissions. Emails will be addressed within 3 business days, including a detailed plan to investigate the issue and any potential workarounds to perform in the meantime. Do not report non-security-impacting bugs through this channel. Use [GitHub issues](https://github.com/goharbor/harbor/issues/new/choose) instead.

### Proposed Email Content
Provide a descriptive subject line and in the body of the email include the following information:
* Basic identity information, such as your name and your affiliation or company.
* Detailed steps to reproduce the vulnerability  (POC scripts, screenshots, and compressed packet captures are all helpful to us).
* Description of the effects of the vulnerability on Harbor and the related hardware and software configurations, so that the Harbor Security Team can reproduce it.
* How the vulnerability affects Harbor usage and an estimation of the attack surface, if there is one.
* List other projects or dependencies that were used in conjunction with Harbor to produce the vulnerability.
 
## When to report a vulnerability
* When you think Harbor has a potential security vulnerability
* When you suspect a potential vulnerability but you are unsure that it impacts Harbor
* When you know of or suspect a potential vulnerability on another project that is used by Harbor. For example Harbor has a dependency on Docker, PGSql, Redis, Notary, Clair, etc.
  
## Patch, Release, and Disclosure
The Harbor Security Team will respond to vulnerability reports as follows:
 
1.  The Security Team will investigate the vulnerability and determine its effects and criticality.
2.  If the issue is not deemed to be a vulnerability, the Security Team will follow up with a detailed reason for rejection.
3.  The Security Team will initiate a conversation with the reporter within 3 business days
4.  If a vulnerability is acknowledged and the timeline for a fix is determined, the Security Team will work on a plan to communicate with the appropriate community (to be completed within 1-10 days of the report of the vulnerability), including mitigating steps that affected users can take to protect themselves until the fix is rolled out.
5.  The Security Team will also create a [CVSS](https://www.first.org/cvss/specification-document) using the [CVSS Calculator](https://www.first.org/cvss/calculator/3.0). The Security Team makes the final call on the calculated CVSS; it is better to move quickly than making the CVSS perfect. Issues may also be reported to [Mitre](https://cve.mitre.org/) using this [scoring calculator](https://nvd.nist.gov/vuln-metrics/cvss/v3-calculator).
6.  The Security Team will work on fixing the vulnerability and perform internal testing before preparing to roll out the fix.
7.  The Security Team will first email the fix to cncf-harbor-distributors-announce@lists.cncf.io, so that they can further test the fix and gather feedback. See the section **Disclosure to Private Distributors List** for details about how to join this mailing list.
8.  Once the fix is confirmed, the Security Team will patch the vulnerability in the next patch or minor release, and backport a patch release into all earlier supported releases. 

### Fix Disclosure Process
The Security Team publishes an [advisory](https://github.com/goharbor/harbor/security/advisories) to the Harbor community via GitHub. In most cases, additional communication via Slack, Twitter, CNCF lists, blog and other channels will assist in educating Harbor users and rolling out the patched release to affected users. 

The Security Team will also publish any mitigating steps users can take until the fix can be applied to their Harbor instances.
 
## Mailing lists
- Use cncf-harbor-security@lists.cncf.io to report security concerns. Received by Product Security Team members, and used by the security team to discuss security issues and fixes.
- Join cncf-harbor-distributors-announce@lists.cncf.io for early private information, including mitigating steps, on security patch releases. See below for information on how Harbor distributors can apply to join this list.

## Disclosure to Harbor Distributors List
This private list is intended to be used primarily to provide actionable information to multiple distributor projects at once. This list is not intended to inform individuals about security issues.

### Membership Criteria
To be eligible to join the cncf-harbor-distributors-announce@lists.cncf.io mailing list, you should:
1. Be an active distributor of Harbor.
2. Have a user base that is not limited to your own organization.
3. Have a publicly verifiable track record up to the present day of fixing security issues.
4. Not be a downstream or rebuild of another distributor.
5. Be a participant and active contributor in the Harbor community.
6. Accept the Embargo Policy that is outlined below. 
7. Have someone who is already on the list vouch for the person requesting membership on behalf of your distribution.

**The terms and conditions of the Embargo Policy apply to all members of this mailing list. A request for membership represents your acceptance to the terms and conditions of the Embargo Policy**

### Embargo Policy
The information that members receive on cncf-harbor-distributors-announce@lists.cncf.io must not be made public, shared, or even hinted at anywhere beyond those who need to know within your specific team, unless you receive explicit approval to do so from the Harbor Security Team. This remains true until the public disclosure date/time agreed upon by the list. Members of the list and others cannot use the information for any reason other than to get the issue fixed for your respective distribution's users.
Before you share any information from the list with members of your team who are required to fix the issue, these team members must agree to the same terms, and only be provided with information on a need-to-know basis.

In the unfortunate event that you share information beyond what is permitted by this policy, you must urgently inform the cncf-harbor-security@lists.cncf.io mailing list of exactly what information was leaked and to whom. If you continue to leak information and break the policy outlined here, you will be permanently removed from the list.
 
### Requesting to Join
Send new membership requests to cncf-harbor-security@lists.cncf.io.
In the body of your request please specify how you qualify for membership and fulfill each criterion listed in the Membership Criteria section above.

## Confidentiality, integrity and availability
We consider vulnerabilities leading to the compromise of data confidentiality, elevation of privilege, or integrity to be our highest priority concerns. Availability, in particular in areas relating to DoS and resource exhaustion, is also a serious security concern. The Harbor Security Team takes all vulnerabilities, potential vulnerabilities, and suspected vulnerabilities seriously and will investigate them in an urgent and expeditious manner.

Note that we do not currently consider the default settings for Harbor to be secure-by-default. It is necessary for operators to explicitly configure settings, role based access control, and other resource related features in Harbor to provide a hardened Harbor environment. We will not act on any security disclosure that relates to a lack of safe defaults. Over time, we will work towards improved safe-by-default configuration, taking into account backwards compatibility.
