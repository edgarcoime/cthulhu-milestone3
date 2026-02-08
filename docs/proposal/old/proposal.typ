#let main_color = rgb("#1E2124")
#let main_gray = rgb("#232837")
#let main_font = "Noto Sans"

#set document(
  author: "Edgar Jaycinth Coime",
  title: "Cthulhu Project Proposal",
  description: "Cthulhu Project Proposal",
  date: auto
)
#set page(paper: "a4")
#set text(fill: main_color)

#let title-page(title:[], sub: [], email:[], osid:[], name: "Edgar J Coime", tp: true, fill: white, body) = {
  set page(fill: rgb(fill), margin: (top: 1.5in, rest: 2in))
  set text(font: main_font, size: 12pt, weight: "regular")
  align(left)[
    British Columbia Institute of Technology 
    #v(0.1em)
    #text(size: 10pt, "Network Security Applications Development")
  ]
  line(start: (0%, 0%), end: (8.5in, 0%), stroke: (thickness: 4pt))
  align(horizon + left)[
    #text(size: 24pt, title, weight: "medium")\
    #v(1em)
    #text(sub)
    #v(2em)
    #name \
    #osid \
    #link(email)
  ]
  
  align(bottom + left)[#datetime.today().display()]
  
  if tp {
    pagebreak()
    set page(fill: none, margin: auto)
    align(top, outline(indent: auto))
  }
  pagebreak()
  body
}

#show: body => title-page(
  title: [Project Proposal \ Cthulhu],
  sub: [COMP 8800 Major Project \ Ashkan Jangodaz],
  email: "mailto:ecoime1@my.bcit.ca",
  osid: [A01003601],
  fill: white,
  body
)

// AFTER TABLE OF CONTENTS MAIN CONTENT SETTINGS
#set page(
  fill: white, 
  margin: (top: 1.25in, bottom: 1.25in, rest: 1.0in), 
  numbering: "1 of 1",

  // Footer
  footer: context [
    Edgar J Coime
    #h(1fr)
    #counter(page).display(
      "1 of 1",
      both: true,
    )
  ]
)
#set heading(numbering: "1.1.1")
#show heading: set block(below: 1em)
#set text(font: main_font, size: 11pt)
#set par(
  first-line-indent: (amount: 0.5in, all: true),
  leading: 1.0em,
)

= Student Background
My name is Edgar Coime, and I am in my 3rd term of the Network Security Applications option of the Applied Computer Science program in BCIT. As per my option, my area of expertise is mainly on Network Security but throughout my tenure at BCIT I have had plenty of experience creating Full Stack Web Applications. One of the projects that I am most proud of is creating an AI learning assistant in 2022 for the Stormhacks Hackathon where two of my schoolmates and I won first place for our inventiveness and overall usability and appeal towards majority of the audience and judges.

I have also had the pleasure of working in the industry for my COOP term in the Computer Systems Technology program in BCIT. During that work term I served the role as project lead helping guide the project through its desired deliverables and communicating directly with the client translating his requirements into user stories and helping his vision come to life. I was also responsible for onboarding future teams into the project so that they can fully understand the codebase and structure of the live application.

= Project Summary 
CTHULHU is an Anonymous-first file sharing platform designed to address privacy and usability shortcomings in modern file sharing systems. The application will enable users to upload and share files/folders of up to 1GB with automatic deletion after 48 hours for anonymous users. Authenticated users will have more fine grain control to extend file lifetimes up to 14 days. Built using modern microservice architecture, CTHULHU will use React NEXT for the frontend and Golang for backend services, it will use HTTP as the communication layer between the frontend web application and API gateway, while taking advantage of RabbitMQ message queuing for inter-service communication. This architecture ensures scalability, maintainability, and robust security while providing a seamless user experience without traditional account barriers.

= Problem Statement and Background
In today’s modern age file sharing has advanced into more than just a convenience and has become an essential tool that many teams across every industry and even students rely on. My fellow students and I often have trouble finding the easiest and most efficient way to share large project material with each other. Current file sharing platforms such as Microsoft OneDrive and Google Drive require a complex flow of cumbersome protocols that increase friction when a user only wants to share a file. They force you into mandatory account creation, multi-step authentication, and complex sharing workflows that involve email-based invitations. These steps create unnecessary friction and prioritize user profiling over spontaneous file sharing and accessibility.

Beyond these usability issues, existing platforms lack a security-first design philosophy. Most mainstream services prioritize data collection and persistent storage over user privacy, creating vulnerabilities that have led to data breaches that have affected millions of users all around the globe. These platforms often store files indefinitely, excessively track their users, and offer limited anonymity protections which create an environment where convenience comes at a cost of user protection and privacy.

= Inventiveness
CTHULHU’s innovation is in its focus on temporal file management and anonymity-fist architecture that eliminate traditional barriers while maintain security integrity. Unlike existing platforms that require persistent user accounts and complex authentication workflows, CTHULHU enables immediate file sharing though its account-optional functionality with a fix 48-hour time bound and its assurance of automatic deletion when that time has elapsed. Unlike more contemporary platforms account creation is not mandatory and only enable the user access to more fine grain control over their files.

The system will also be Integrating advanced microservice architecture utilizing RabbitMQ message queuing for inter-service communication. Microservices is an exciting paradigm and will enable the platform to grow according to its needs. It allows us to add additional features (i.e.. Logging service, notification service, virus scanning service, etc.) without radically changing the codebase, which Monolithic architectures would demand giventhe same change.

Test @Kumar2023Implementation


= Complexity
Building CTHULHU will involve intensive research on microservice architecture and overcoming the challenges that I have never encountered in any of my previous projects. CTHULHU will demand an understanding of distributed systems and how to effectively manage the asynchronous inter-service communication that comes from RabbitMQ’s message queue system. It would require implementing a robust message queue protocol, ensuring data consistency across multiple independent services, and making sure that data remains ACID compliant during transactions that may need to be rolled back if catastrophic failure occurs.Event driven architecture is also something that I have never had experience with and is quite infamous for its difficulty. It would require overcoming complex service orchestration challenges and understanding how to handle transactions through CAP theorem and SAGA patterns.

File upload functionality also introduces critical security vulnerabilities that require mitigation strategies for: preventing XSS attacks through malicious file uploads, implementing file type validation and sanitation, defending against path traversal attacks, and many more. The technical complexity starts to balloon when we account for the sophisticated message broker architecture, distributed system fault tolerance, and malware detection integration. Therefore the complexity for CTHULHU I feel is more than sufficient for this major project.

= Technical Challenges
Most of the applications I have built have been primarily monolithic architectures with the occasional mix of microservice through client and server. CTHULHU gives me the opportunity to fully dive into the inner workings of microservice architecture applications especially when incorporating an asynchronous communication layer through RabbitMQ. Through this project I will have to learn about microservice methodologies and communication patterns such as event stream pub/sub, event sourcing, Command Query Responsibility Segregation (CQRS), and the SAGA pattern for distributed transaction management.

Service-to-Service authentication and Authorization in microservices is foreign to me and is completely different to more synchronous communication patterns like HTTP. In HTTP we can make use of existing user sessions but microservices require machine-to-machine communication usually through JWT tokens or mutual TLS. But we also need strict access control policies to mitigate risks from compromised services or malicious lateral movement in the network.

Testing references @Atovi2022Microservice, 

// Citation Section
#pagebreak()
#bibliography("references.bib", style: "american-psychological-association", title: "References")
