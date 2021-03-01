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
  - [X] Rebase and deploy
  - [x] Provide good arguments for choice of virtualization techniques and deployment targets
- [X] Log dependencies

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

**Further notes:** We moved from Dockerhub to Azure Container Registry as this seems to be faster with Azure App Service and has better support for it.

## Switching to PostgreSQL

As we are already using the Microsoft Azure ecosystem, we decided it would be the ideal place to also host our database. We wanted to still use a relational database (our models have clear relations between them), but we decided to use another relational database - Postgres, as it has support for more data types, which would give us more benefits in the future. We also considered the case that if we have to do scale the application on a database layer, SQLite has limited concurrency. As another advantage, Azure offers hosting and scailing of PostgreSQL databases, which was another reason for our change of database.

## Week 04 Continuous Integration (CI), Continuous Delivery (CD), and Continuous Deployment

- [X] Complete implementing an API for the simulator
- [X] Creating a CI/CD setup for your ITU-MiniTwit.

### Choice of CI/CD provider

We have chosen to go with Azure DevOps Pipelines as our CI/CD provider. The reasoning behind this was, in large part, that we were already using the platform for our sprint backlogs. In addition, it also provides good integration with our platform for deployments i.e. Azure, which is our main reason for this choice.

We are running a CI pipeline on our develop branch in order to verify that the code we are continuously contributing to the project does not break any of the existing codebase. In our CI pipeline we are ensuring that all test are still passing, and the program is able to be compiled.

**CI Piplines:**
- CI Test: is a test pipeline that builds and execute all unit tests. This pipeline is triggered with every PR and has to pass before a merge can happen (configured as branch rules in Github).
- CI build: builds a dockrt image and pushses it to Azure Container Registry (ACR) and publishes it as an artifact we can use in the CD pipeline. This pipeline is triggered whenever a PR is merged into main.

**CD Pipeline / Release:**
- Deploy to App Service: is configured to take the artifact from the CI build pipeline, which is the latest docker image from ACR, and execute a series of tasks. The first task called "Deploy Azure App Service" specifies a connection to our ACR and Azure App Service and tells the App Service to pull the latest image. The second task simply restarts our app service. The third tasks creates a release on Github where it increments the release tag, attach assets from the artifact and add a changelog based on the commit history. This pipeline is triggered whenever a new image is pushed to ACR, thus the whole pipeline will be executed whenever new code is merged into our main branch. This ensures Continous deployment.

**Further notes:** As of right now we are considering our project Open Source, meaning that some measures described in this pipeline is based on other people contributing to the project. Was this to be closed source we would consider a Trunk based branching strategy with less manual steps.

### New project structure

Since part of next weeks work will be "cleaning and polishing of our ITU-MiniTwit" application, we decided upon a new project structure for our application that we will be implementing by then. The overall goal of this refactoring will be to make the code more readable, maintainable, and easier to deploy.

The reason this refactoring is necessary so soon is, out initial refactoring from python to golang was, very literally, a 1-1 translation from the python application. This has resulted in our current application having no separation in responsibilities in regard to which class does what, as well as the UI, and the API being to separate applications that need to be deployed. Since the API and the UI is each contained entirely in their own class, there is a lot of code duplication as well between the two.

Our current idea is to follow the overall structure proposed [here](https://github.com/Mindinventory/Golang-Project-Structure). The API will therefore be merged into UI i.e. there will only be one application. This part of the work was already started this week. The data access layer will also be split into different services from the http handlers. The ending project structure should end up looking like the following, with the exception that we will only have a single version of the api that we will be maintaining. 

![](https://raw.githubusercontent.com/Mindinventory/Golang-Project-Structure/master/structure.png)

## Week 05 Finalizing development stage

- [ ] Add API tests to CI
- [ ] DevOps - the "Three Ways"
- [x] Software Maintenance || - Group B

### Group B Monitoring

We performed the following tests on Group B's minitwit app, hosted on this [link](http://144.126.244.138:4567)

- Do you see a public timeline? - Yes
- Does the public timeline show messages that the application received from the simulator? - Waiting on simulator
- Can you create a new user? - Yes, [example](http://144.126.244.138:4567/petya)
- Can you login as a new user? - Yes, with the same user as above
- Can you write a message? - Yes, messages can be found [here](http://144.126.244.138:4567/public) and [here](http://144.126.244.138:4567/petya)
- After publishing a message, does it appear on you private timeline? - Yes, messages can be found [here](http://144.126.244.138:4567/petya)
- Can you follow another user? - Yes, however we experience a bug that on a freshly created user, there are few automatically assigned accounts that the user follows. - We reported the issue to group B [here](https://github.com/DevOps2021-gb/devops2021/issues/44)
