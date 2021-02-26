# Weekly Log

## Week 01

- [x] Add version control
  - [Repo](https://github.com/heyjoakim/devops-21)
- [x] Try to develop a high-level understanding of ITU-MiniTwit.
- [x] Migrate ITU-MiniTwit to run on a modern computer running Linux

  - [x] Get python to run
  - [x] Install deps
  - [x] Recompile flag_tool
  - [x] Install SQLite browser
  - [x] Run 2to3 to convert py2 to py3
  - [x] `shellcheck` and fix `control.sh`

- [x] Share Work on GitHub
  - [Repo](https://github.com/heyjoakim/devops-21)
- [x] Prep for next week
  - [x] Discussed branching strategy, explained below in the notes section

## Notes

We meet on Mondays from 10.00 - X.X.X.X (Super agile here!!)

### This is our branching strategy

PR
Develop -> Feature

Branch out from develop into a feature / bug and then create a PR to merge back into develop. From develop releases are pushed to production (maybe one test environment?).

## Week 02

- [x] Choose language and technology for refactoring
  - [x] And why
- [x] Choose branching strategy
- [x] Refactor
- [x] Commitment guidelines?
- [x] Implement API for simulator

### Choose language and technology

| Lang/Dev                | Pros                                                                                                            | Cons                                                                                                                                                          |
| ----------------------- | --------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Go/Gorilla              | Fast compared to other suggested frameworks [[1]](https://github.com/the-benchmarker/web-frameworks), fullstack | Setting up env can be tricky                                                                                                                                  |
| C#/ASP.NET/Razor/Blazor | Scalable, plenty of resources, fullstack                                                                        | Somewhat heavy framework [[1]](https://github.com/the-benchmarker/web-frameworks), not easy to make 1:1 mapping, due to different structure, Still early life |
| JS/Angualr              | Strong community, fullstack                                                                                     | Not suited for 1:1 app, Not statically type                                                                                                                   |
| JS/Vue                  | Easy to get started with, lightweight                                                                           | Not suited for 1:1 app, Not statically type, Needs separate backend                                                                                           |

We have chosen Go as we believe this is well suited for such task and is fast compared to other frameworks.

### Choose a branching strategy

We are discussing advantages and disadvantages between a Git Flow and Topic/Feature workflow strategy.

| Strategy         | Pros                                                                 | Cons                             |
| ---------------- | -------------------------------------------------------------------- | -------------------------------- |
| Git flow         | Separate releases, more "controlled", more suited for weekly release | "more" work                      |
| Feature workflow | Continous development, cleaner Git history, Simple, Faster deploys   | Need more internal communication |

We have chosen to go with a modified Git Flow strategy as we believe this is more suited for our weekly releases. We have decided to omit the release branch, since we think that it would create unnecessary overhead compared to the size of the project. Our _development_-branch will do tests once CI/CD i setup.

![](https://i.imgur.com/ea6o39W.png)

The branch structure will therefore be as following :

- `develop` All new feature branches must check out from here into feature branches and merged back into develop. The contents of the development branch would usually reflect what is deployed to the test environment.
- `main` The production branch reflects the current deployment in production. The production branch is merged with the develop branch every time a new version deployed to production.
- `feature/{feature-name}` New features are developed on feature branches following the _feature / feature name branch_ structure.
- `hotfix\{hotfix-name}` New hot fixes are developed on separate hot fix branches following the _hotfix / hotfix branch name_

## Week 03 Virtalization

- [x] Complete implementing an API for the simulator
- [x] Continue refactoring
  - [x] Introduce a DB abstraction layer
  - [x] Arguments for choice of ORM framework and chosen DBMS
  - [ ] Rebase and deploy
  - [x] Provide good arguments for choice of virtualization techniques and deployment targets
- [ ] Log dependencies

#### Release and deploy

Azure as cloud provider with docker!

#### ORM Tool

We decided to use [GORM](https://github.com/go-gorm/gorm) as it is one of the most widely used ORMs for Golang (reference:https://github.com/go-gorm/gorm) and also after further research we found it to be the most well-documented.
We also discussed switching to PostgreSQL as a datasource, but decided to postpone that to a later stage, as the ORM abstraction will give us the flexibility to change data storages.

#### Choice of ORM

So far, the application had been constructing it's own SQL statements, and executing them as prepared statements, using SQLite3. However, we need to find a way to best prepare ourselves and minitwit for any changes that may have to be done.

Gorm makes it possible for us to use the golang structs that we already have been working with, in such a way that we can save our objects directly to the database, thereby also having a more explicit struct-strategy in our code.
We also expect that later in the course, it might become necessary to do some refractoring of the database, which is easier with the code-first workflow of Gorm. In that respect, we expect to be able to more dynamically manipulate our database codefirst. Creating primary keys, columns and rows can all be manipulated and created code first.
We also hope to be able to get rid of some repetetive boilerplate SQL, and thus make the code more readable, to the non SQL initiated developer.

Another positive benefit could be that changing to another dbms, could require less work in terms of rewriting code, thus improving modifiability.

#### Choice of virtualization techniques and deployment targets

For hosting our Minitwit App and API, we decided to use Microsoft Azure as a cloud provider. Azure also allows to deploy an application in a Docker containter instance, which was the initial reason why we preferred Docker as a virtualization technique. Moreover, the team had some previous experience and general prefference towards using Docker. Therefore, we created a Docker image, published it on DockerHub and created the basis for further optimizations in our CI/CD pipeline.

## Switching to PostgreSQL

As we are already using the Microsoft Azure ecosystem, we decided it would be the ideal place to also host our database. We wanted to still use a relational database (our models have clear relations between them), but we decided to use another relational database - Postgres, as it has support for more data types, which would give us more benefits in the future. We also considered the case that if we have to do scale the application on a database layer, SQLite has limited concurrency. As another advantage, Azure offers hosting and scailing of PostgreSQL databases, which was another reason for our change of database.

## Week 04 Continuous Integration (CI), Continuous Delivery (CD), and Continuous Deployment

- [] Complete implementing an API for the simulator
- [] Creating a CI/CD setup for your ITU-MiniTwit.
