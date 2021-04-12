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
  - [x] Rebase and deploy
  - [x] Provide good arguments for choice of virtualization techniques and deployment targets
- [x] Log dependencies

#### Release and deploy

Azure as cloud provider with Docker!

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

For hosting our Minitwit App and API, we decided to use Microsoft Azure as a cloud provider because Azure has a solid Docker integration, compared to other cloud providers. Azure also allows to deploy an application in a Docker containter instance, which was the initial reason why we preferred Docker as a virtualization technique. Moreover, the team had some previous experience and general prefference towards using Docker. Therefore, we created a Docker image, published it on DockerHub and created the basis for further optimizations in our CI/CD pipeline.

**Further notes:** We moved from Dockerhub to Azure Container Registry as this seems to be faster with Azure App Service and has better support for it.

## Switching to PostgreSQL

As we are already using the Microsoft Azure ecosystem, we decided it would be the ideal place to also host our database. We wanted to still use a relational database (our models have clear relations between them), but we decided to use another relational database - Postgres, as it has support for more data types, which would give us more benefits in the future. We also considered the case that if we have to do scale the application on a database layer, SQLite has limited concurrency. As another advantage, Azure offers hosting and scailing of PostgreSQL databases, which was another reason for our change of database.

## Week 04 Continuous Integration (CI), Continuous Delivery (CD), and Continuous Deployment

- [x] Complete implementing an API for the simulator
- [x] Creating a CI/CD setup for your ITU-MiniTwit.

### Choice of CI/CD provider

We have chosen to go with Azure DevOps Pipelines as our CI/CD provider. The reasoning behind this was, in large part, that we were already using the platform for our sprint backlogs. In addition, it also provides good integration with our platform for deployments i.e. Azure, which is our main reason for this choice.

We are running a CI pipeline on our develop branch in order to verify that the code we are continuously contributing to the project does not break any of the existing codebase. In our CI pipeline we are ensuring that all test are still passing, and the program is able to be compiled.

**CI Piplines:**

- CI Test: is a test pipeline that builds and execute all unit tests. This pipeline is triggered with every PR and has to pass before a merge can happen (configured as branch rules in Github).
- CI build: builds a docker image and pushes it to Azure Container Registry (ACR) and publishes it as an artifact we can use in the CD pipeline. This pipeline is triggered whenever a PR is merged into main.

**CD Pipeline / Release:**

- Deploy to App Service: is configured to take the artifact from the CI build pipeline, which is the latest docker image from ACR, and execute a series of tasks. The first task called "Deploy Azure App Service" specifies a connection to our ACR and Azure App Service and tells the App Service to pull the latest image. The second task simply restarts our app service. The third tasks creates a release on Github where it increments the release tag, attach assets from the artifact and add a changelog based on the commit history. This pipeline is triggered whenever a new image is pushed to ACR, thus the whole pipeline will be executed whenever new code is merged into our main branch. This ensures Continous deployment.

**Further notes:** As of right now we are considering our project Open Source, meaning that some measures described in this pipeline is based on other people contributing to the project. Was this to be closed source we would consider a Trunk based branching strategy with less manual steps.

### New project structure

Since part of next weeks work will be "cleaning and polishing of our ITU-MiniTwit" application, we decided upon a new project structure for our application that we will be implementing by then. The overall goal of this refactoring will be to make the code more readable, maintainable, and easier to deploy.

The reason this refactoring is necessary so soon is, our initial refactoring from Python to Go was, very literally, a 1-1 translation from the python application. This has resulted in our current application having no separation in responsibilities in regard to which class does what, as well as the UI, and the API being to separate applications that need to be deployed. Since the API and the UI is each contained entirely in their own class, there is a lot of code duplication as well between the two.

Our current idea is to follow the overall structure proposed [here](https://github.com/Mindinventory/Golang-Project-Structure). The API will therefore be merged into UI i.e. there will only be one application. This part of the work was already started this week. The data access layer will also be split into different services from the http handlers. The ending project structure should end up looking like the following, with the exception that we will only have a single version of the api that we will be maintaining.

![](https://raw.githubusercontent.com/Mindinventory/Golang-Project-Structure/master/structure.png)

## Week 05 Finalizing development stage

- [ ] Add API tests to CI
- [x] DevOps - the "Three Ways"
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

changeset 0.1

### The Three Ways

Some methodologies from the _Three Ways_ won't be covered, since they don't relate to how we work on this project.

**Flow**

_Make work visible_

To make our scheudled work visible throughout the team, we make use of a _taskboard_ in Azure DevOps. All tasks are children to a _user story_ that describes a goal or issue on a higher level. Specific tasks are then created under each user story. The entire team can inspect this board and get an overview on how far the team is with the project work and who works on what.

_Limit Work in Progress & Reduced batch sizes_

When creating tasks in the team, we ensure to make the tasks small enough to be addressed withing a relative short period of time. We strive after not having tasks that span over mulitple days. This also makes the merging way more simple.

_Reduced number of handoffs_

To reduce the numberof handoffs and delaying steps, like having one person to do code review and accept pull requests, the entire team have these capabilities. This ensures that once a pull request are created, the first available collegue can review it start the automatic testing process in out CI pipeline. The same applies to deployment.

**Feedback**

Creatig fast feedback is critical to ensure quality and reliability.

_See problems as they occur_

To catch any errors in the code, we make use of automated testing in our CI pipeline. The pipeline runs a series of tests that ensures functionality has not been broken by the new additions. If a test fails, the pipeline terminates and the team members will get an email that tells that tests has failed. This is then up to the comitter to fix any issues.

_Swarm and solve problems_

In relation to the above, our branching strategy allows us to commit to a failing branch. If we discover an error in the production application, we are allowed to create a hotfix branch from _main_, and create a pull request directly into main once the issue is addressed.

_Pushing quality closer to source_

We strive the distribute the reposnsibility of certain parts of the system across the team members, so one person is the expert, but everyone in the team should know how all components interact. The _expert_ should be able to explain and handover tasks to other team mates with confidence though.

**Continual Learning and Experimentation**

_Enabling organizational learning and a safety culture_

Our project builds upon the generative culture, which means that we don't blame people for any mistakes, but sees it as a joint learning experience. When on team members discovers an issue and asks for help, most often one person knows the answer and shares it with the entire team. This ensures that errors and mistakes only happen once, which generally leads to more work on new features, instead of struggling with the system.

### Software Maintenance

We have not encountered nor received new bug reports

## Week 06 Monitoring

- [x] Find metrics and user statistics for system (prep material)
- [x] Add monitoring
- [x] Software Maintenance

### Metrics and user statistics

#### CPU load during the last hour/the last day

**Stakeholder:** Developer

In Azure Portal, it is straight forward to find CPU usage of an _App Service_. Go to **Diagnose and solve problems** > **Availability and Performance** > **CPU Usage**. The resolution of CPU usage data is 5 minute intervals. The average CPU usage between 23.00 and 23.59 the 6th of March was **0.18%**.

#### Average response time of your application's front page

**Stakeholder:** Developer

We can't get a reading for the front page, but we can get data from the entire site. Go to **Monitoring** > **Metrics**, then select **minitwut,Response Time,Avg** and select a time range and granularity. The average reponse time for the sustem the last 3 days was **28.88ms**.

#### Amount of users registered in your system

**Stakeholder:** C-level officers

We can query our database to get the information. The following query returns the number of users in the system:

```sql
SELECT COUNT(*) from "user";
```

The number of users registered in the system: **8926** (2020-03-07 20.28).

#### Average amount of followers a user has

**Stakeholder:** C-level officers

We can query our database to get the information. The following query returns the average number of followers each user has:

```sql
SELECT AVG(flws) FROM (SELECT COUNT(*) AS flws FROM "user" LEFT JOIN follower ON (follower.who_id = "user".user_id) GROUP BY "user".user_id) _;
```

The average number of followers a user has: **2.01 ~ 2** (2020-03-07 20.53).

### Monitoring with Prometheus and Grafana

In week 07 we created a Prometheus server, hosted on Digital ocean: `http://142.93.103.26:9090`, which contains the metrics tracked by prometheus and custom metrics for CPU and memory usage, total numbers of users and messages.

All of the metrics can be found also on: `https://minitwut.azurewebsites.net/metrics`, the custom metrics are defined as: `group_l_minitwut_*`

Once we setup the metrics tool, we used Grafana to consume and visualize the Prometheus data. We created a dockerfile for Grafana, downloading the latest image and specifying a custom data source and dashboard. All the settings can be found in the grafana folder. We then deployed the dashboard to Digital ocean: `http://164.90.165.111:3000/d/JJQvP88Mz/prometheus-2-0-stats`. By defining the Grafana configuration, we avoid loosing our metrics everytime when we restart the container.

We decided to explore a different cloud provider - Digital ocean, as we encountered limitation on the free-tier plan in Azure, where our minitwit app is running. Before making the decision to fully migrate our minitwit app, we decided to use the monitoring tools as a try-out of another provider.

### Software Maintenance

We have not encountered nor received new bug reports

## Week 07 Monitoring

- [x] Add Maintainability and Technical Debt estimation tools to your projects
- [x] Enhance your CI Pipelines with at least two static analysis tools
- [x] Software Maintenance

### Add Maintainability and Technical Debt estimation tools to your projects

We integrated to our project the following tools:

- [x] Sonarqube: 3 bugs, 0 vulnerabilities, 16 code smells
- [x] Code Climate: 3 days maintainability, Technical debt is in steady decline since the beginning of the project (from almost 11% to 3.9%)
- [x] Better Code Hub

We found several issues relating to old Python code that we still kept in the repository and the api simulation tests. As well as some code smells. We solved those issues by removing old code, fixing bugs, rewriting the Python tests in Golang and refactoring repetitive code.

### Enhance your CI Pipelines with at least two static analysis tools

### Software Maintenance

We found some issues reported by the static analysis tools and fixed them accordingly.

### New brahnching strategy

Previously we noted the use of a Git Flow branching strategy, meaning that contributions were to be made off of a branch from develop, though due to the addition of Static Code analysis and the rewrite of API simulation tests from Python to Golang, we have decided to adopt a new strategy. The motivation for the change in stratgy is a move towards Continous deployment. Our pipelines now support a full deployment process from main branch to our productions servers with automatic releases. This means that a Git Flow strategy will only hinder our lead time from user story to production. As this project is a Open Source project we do not intend to adopt a trunk based stragy but a Github Flow strategy.

##### Github Flow branching strategy

The adoption of this new strategy means that branches are to be made from the `main` branch, with the one rule that `main` always should be deployable. When changes are made you open a pull request whereas a discussion and review of the code will be conducted.

## Week 08 Logging

- [x] Add Logging to Your Systems
- [x] Test that Your Logging Works
- [ ] Write an SLA for Your Own API

### Add Logging to Your Systems

For logging, we integrated [Logrus](https://github.com/sirupsen/logrus), as it supported various formats for the logs, and it's something we deemed necessary in order to integrate them into a logging system.

We had some problems with the EFK stack and started exploring other possiblities, such as DataDog.

### Test that Your Logging Works

When this is being written, a member has introduced an error somewhere in the application, and now another member will investigate by viewing the logs.

Looking at the logs, this error shows up:

![Screenshot from 2021-04-04 17-40-24](https://user-images.githubusercontent.com/43805989/113514120-21719d80-956d-11eb-9cf0-ae5304a4107a.png)

```

{ err crypto/bcrypt: hashedPassword is not the hash of the given password
  level error
  time 2021-04-04T13:30:58Z
}
```

This looks rather strange. The errors occurs in conjunction with a error-level log regarding log-in. Therefore I assume that this hashing error occurs in connection with login.

Now, I'll attempt to recreate the error in production. For doing this, I create a user with the following credentials:

- username : Testuser
- email : usr@email.com
- Password : pw

And then attempt to log in with these credentials. This throws the same error and I cannot login. The message "Invalid Password" is shown.

I start looking around in the controller for logins. In line 42 I find this code, where the digest of the incoming password is compared to the stored password.

```
if err := bcrypt.CompareHashAndPassword([]byte(user.PwHash + "a"), []byte(r.FormValue("user"))); err != nil {

```

An "a" is appended to the end of the digest from the database. This seems strange. I try removing the "a".

Trying to test this in localhost, the problem still occurs.
I then notice that this error is reading the "user" field from the formvalue, I change it to "password".

It now works nicely as it should.

### Write an SLA for Your Own API

To be completed

## Week 09 Scailing

- [ ] Add Scaling to your projects
- [x] Software Licensing
- [x] Software Maintenance
- [x] Logging with DataDog
- [x] Extra: Migrating from Azure to Digital Ocean

### Adding a license file

Regarding the choice of a license file, we chose the MIT license. This license seems to be the most widely used open-source license, which is also reflected in our dependencies. With the MIT license, we are not requiring that people modifying our application also are open-sourced.

Having this license, we also ensure that we comply with any gpl-licenses of our dependencies that would require us to be be open source in any capacity.

### Logging with DataDog

Due to issue with the ELK stack and our current setup with Microsoft Azure, which was imposing limitations due to the inability (or our lack of knowledge) to open new ports and have more control over the server setup, we decided to use DataDog.

### Migrating from Azure to Digital Ocean

At the beginning of this project we argued our choice for Microsoft Azure, however further into the project we ended up in a vendor lock-in, where our Docker image was hosted in the Azure Container Registry, our pipeline was in Azure DevOps, our PostgreSQL database was hosted in Azure. Once we ran out of the free-tier student credits and we realised we need to switch to another cloud provider that does not charge us so much for the same services and we can use more student credits, it became hard for us to migrate to another provider. We have had prior experience with Digital ocean for our monitoring when we initially started considering switching providers. Based on the experience we have already acquired and the exchange with the other Masters groups in the course, Digital Ocean was the safe choice to switch to.

When switching we considered a setup that will not impede us in the future from moving to another provider:

- We had to change the publishing of our Docker image from Azure Container Registry to DockerHub, this way - regardless of the provider, we can always pull our iamge free of charge
- We switched from Azure DevOps to Github Actions, thus allowing us to swap the target server to deploy to, and our pipeline will not be broken
- We also migrated our Postgres database do digital ocean, as the hosting in Azure was becoming quite expensive. The database is hosted seperately, so it is not tied to a service or a provider - can be easily imported in another storage provider.

Inevitably, we lost data in the migration process. And the logging system we setup last week was very useful when we encountered migration issues and we could quickly find the problem in our system.

App url: http://206.189.14.172:8000/public
Server url: http://206.189.14.172:8000/api

## Week 10 Scailing

- [ ] Add Scaling to your projects
- [ ] Software Licensing
- [x] Software Maintenance

### Adding scailing with XXX

### Software Licensing
