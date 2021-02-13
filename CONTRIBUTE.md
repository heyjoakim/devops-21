# Contributions

We love open source and would happily invite you to contribute to this project.

#### How to contribute to this project
All changes, submissions or suggestions happens through github pull requests. Pleasre read below section about our branching model.

1. Make a fork of the repository and create a branch from `main`
2. Add your changes, features etc. but only 1 feature pr. PR 
3. State explicitly in the PR what your fix or feature does and how to test this.



#### Repository setup
We have a single repository setup with `main`and `develop`branches. At both branches we have applied branch protection rules to require pull request reviews before merging. Thus all commits should be made to non-protected branches e.g. `feature/descriptive_name` and then be submitted through a pull request that requires atleast 1 reviewer before merged. 


#### Branching model

This project is using a modified Git Flow strategy, as we believe this is more suited for our weekly releases. We have decided to omit the release branch, since we think that it would create unecessary overhead compared to the size of the project. Our _development_-branch will do tests once CI/CD i setup. 

![](https://i.imgur.com/ea6o39W.png)

The branch structure will therefore be as following :

- `develop` All new feature branches must check out from here into feature branches and merged back into develop. The contents of the development branch would usually reflect what is deployed to the test environment.
- `main` The production branch reflects the current deployment in production. The production branch is merged with the develop branch every time a new version deployed to production.
- `feature/{feature-name}` New features are developed on feature branches following the *feature / feature name branch* structure.
-  `hotfix\{hotfix-name}` New hot fixes are developed on seperate hot fix branches following the *hotfix / hotfix branch name*


#### Reviewing contributions
Contributions are to be reviewed by atleast 1 direct access contributor. Please do direct any bugs, issues through github and not to the direct access contributers.

Which distributed development workflow will we use?