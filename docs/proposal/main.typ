#import "lab-layout.typ": *
// PAGE CONFIGURATIONS
#let main-font = "Libertinus Serif"
#let doc-title = [Cthulhu Project Proposal]

#set document(
  author: "Edgar Jaycinth Coime",
  title: doc-title,
  description: "Cthulhu Project Proposal",
  date: auto,
)

#show: body => title-page(
  title: doc-title,
  sub: [COMP 8800 Major Project \ Ashkan Jangodaz],
  email: "mailto:ecoime1@my.bcit.ca",
  osid: [A01003601],
  fill: white,
  running-head: "Cthulhu proposal",
  paper-size: "us-letter",
  body
)

= Document Control
== Document Information
#table(
  columns: (auto, 1fr),
  inset: 10pt,
  fill: (x, y) =>
    if y == 0 { gray },
  table.header(
    [], [*Information*]
  ),
  [Document Owner], [Edgar J Coime],
  [Document Owner ID], [A01003601],
  [Issue Date], [2025-10-12],
  [Last Saved Date], [#datetime.today().display()],
  [File Name], [Cthulhu Project Proposal]
)

== Document History
#table(
  columns: (auto, auto, 1fr),
  inset: 10pt,
  fill: (x, y) =>
    if y == 0 { gray },
  table.header(
    [*Version*], [*Issue Date*], [*Changes*]
  ),
  [1.0], [2025-09-19], [Initial Proposal submitted.],
  [1.1], [2025-10-12], [Proposal draft 1.1 submitted reflecting some feedback.],
  [1.2], [2025-10-15], [Added Methodology and Design portion highlighting agile approach.],
  [1.3], [2025-10-18], [Created Test Plan section.],
  [1.4], [2025-10-22], [Created Scope and Depth section detailing application's microservice architecture.],
  [1.6], [2025-10-24], [Added a comprehensive Development Schedule and milestones. Detailing all project deliverables and divided the milestone timeline into even parts for a the 12-week development period.],
  [1.7], [2025-10-26], [Created diagram for System Architecture, File upload flow, and file download flow.],
  [2.0], [2025-10-26], [Added further details in Scope and Depth. Outlined the complexity of Microservice arch. and the depth that comes learning about and mastering RabbitMQ.],
  [2.1], [2025-11-06], [Added Milestone 5: Online Functionality, E2E user tests, and CI/CD pipeline.],
  [2.2], [2025-11-06], [Removed online functionality requirement in milestone 2 and moved it to milestone 5],
  [2.3], [2025-11-08], [Added Milestone 2 details about incorporating cloud storage, how Authentication service should be source of truth, and plan to cache tokens/sessions in gateway.],
  [2.4], [2025-11-09], [Expand Milestone 3 to 4 weeks and add details about designing a file life cycle service database to store persistent data.],
  [2.5], [2025-11-14], [Expand Milestone 4 to 6 weeks. Add details about service performance requirements virus scanning will be computationally expensive. Add SAGA transaction requirements for if file is malicious.],
  [2.6], [2025-11-17], [Added more details for Milestone 4 about virus scanning properties, detection requirements, and test acceptance thresholds.],
  [2.7], [2025-11-21], [Added Milestone 5 details about online functionality and HTTPS thorough valid SSL certs. Add requirement for E2E user testing to ensure Cthulhu platform working properly. Add details about CI/CD pipeline, rollbacks, beta environments, and Docker image registry integration.],
  [3.0], [2025-11-21], [Submission for Proposal draft 3]

)

#pagebreak()
#outline()

#pagebreak()
#include "sections/body.typ"

#pagebreak()
#bibliography(
  "bibliography/ref.bib",
  title: auto,
)
