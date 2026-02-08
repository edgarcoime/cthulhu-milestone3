#import "lab-layout.typ": *
// PAGE CONFIGURATIONS
#let main-font = "Libertinus Serif"
#let doc-title = [Milestone 3 Report]

#set document(
  author: "Edgar Jaycinth Coime",
  title: doc-title,
  description: "Milestone 3 Report",
  date: auto,
)

#show: body => title-page(
  title: doc-title,
  sub: [COMP 8900 Major Project \ Ashkan Jangodaz],
  email: "mailto:ecoime1@my.bcit.ca",
  osid: [A01003601],
  fill: white,
  running-head: "Milestone 3 Report",
  paper-size: "us-letter",
  body
)

#outline()

#pagebreak()
#include "sections/body.typ"
