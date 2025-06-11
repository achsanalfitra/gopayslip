# Deals BE assignment
Study case: payroll system

## Completion

This is my submission. This project is not complete yet; it is about 90% complete. The thing I'm missing is wiring up the routers and handlers

## What is in this repo?

It is a "performance scalable" custom Go framework that is designed to facilitate the business logic. For example, freezing logic is done in-memory and on start-up lookup. Regardless of how many services you spawn, users can't add or modify backdates overtime or reimbursement; this is perfectly valid and convenient in financial/HR settings.

What does this include?
* configurations: server, db
* migrations: sql files
* business logic: all the 7 api endpoints are well defined
* custom auth

I might have underestimated the boilerplate tax when working on this project

## What is the architecture?

When I read the requirement having "performance scalability," this is my Go to. I did have alternatives: out of pockets frameworks like Gin or Chi. Auto orm relation with gorm or goose or even cli tools with soda. Despite that, I chose to write the "bare-metal" framework for Go. I reinvent what is necessary and leave the rest, no bloat, maintainable, and scalable.

The project structure is defined to be modular and semi-monolithic. I've understood that many of the business logic operations are dependent on each other. For example, for proper a audit log send, you have one entry-point that simultaneously send to two tables and the id must be unique. I figured that hooking up multiple microservices introduce 2 things: network overhead and latency. My approach, on the other hand, provides in-memory update, failfast mechanism that allows blocking of invalid methods right at the start.

I think this is architecture cohesive; I have effectively 1 codebase in "internal" to wire up services. Functions, services, templates, and models are reusable and they are within internal. Splitting the services can be as easy as wiring up new config and define the new entrypoint in cmd. For example, one of the admin job is to aggregate the whole payslips to create summary. That can be migrated to cmd/admin/aggregate/main.go if more resource allocation is needed.

## Closing words

I might not be the fastest coder out there and thank you for giving me an opportunity to participate in this test. I've given my best and I hope this doesn't disappoint. This was such a massive lift for my experience. I'll keep improving and try again in future opportunities. 