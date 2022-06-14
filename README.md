# fetch-demo

## What is this?

This is a small monorepo of microservices that look to help keep track of points for users and payers based on the [points exercise](https://fetch-hiring.s3.us-east-1.amazonaws.com/points.pdf).

### Services:

- points service: `port 8500`
- user service: `port 8501`
- payer service: `port 8502`
- swagger-ui: `port 8503`
- nats: ports `4222 6222 8222`

These three services work together to support commands against the points service.  They communicate to one another through Nats Jetstream.

## How do I use this?

Although most setup is automated, there are some dependencies required to get the full stack runnning locally.  Each dependency links to their own installation page.

### Dependencies

- [Go](https://go.dev/doc/install): language code is written in
- [Mage](https://magefile.org/): task automation
- [Docker](https://docs.docker.com/get-docker/): containers :woozy_face:
- [Postman](https://www.postman.com/downloads/): (Optional) api checking

Other tools are used through out this demo for tasks like code generation, debugging, etc.  These are not necessary for running the actual programs.

### Setup

Once you have Go, Mage, and Docker installed, you should be able to start the programs.  In the root of the repository, run:

```bash
mage -v start
```

This command should show some containers starting as well as some outputs from creating streams in Nats.  Everything should be ready to run once this finishes!

### Testing

There is a postman collection saved at the root of the directory.  It contains the configuration for the commands/routes used against the services.  Once this is exported it is pretty straight forward to grab and run the commands.  If you do not wish to use postman, the swagger docs for the services are served on port `8503` so you can see the commands.

#### Order of Operations

In order to test the commands that are listed in the exercise requirements, we first need to add payers and a user to the system.  This can be done through the [addPayer](http://localhost:8503/?urls.primaryName=PointService#/default/addPayer) and [addUser](http://localhost:8503/?urls.primaryName=UserService#/default/addUser) routes.  Ideally, you should add a payer for `DANNON, UNILEVER, and MILLER COORS` to match the test requirements.  The user's name does not matter, however **it's important to grab the user id from the addUser route!**  This will be used for interacting with the points endpoints.

Now that everything is setup, we can go ahead and run the commands in the exercise.  First we need to add all the points to the user through the [givePoints](http://localhost:8503/?urls.primaryName=PointService#/default/givePoints) route.  The id of the user needs to be added to the path in order for this to be successful.  If you hit the [listPoints](http://localhost:8503/?urls.primaryName=PointService#/default/DescribeBalance) route, you should now see your user has points from payers.  YOu can now run the [spendPoints](http://localhost:8503/?urls.primaryName=PointService#/default/spendPoints) command that should return what payers had points spent for the user.  Running the listPoints command again should now show an updated amount of points for the user.


#### Unit Testing

There aren't a ton of tests in this repository, but there are some copied over from older projects as well as some written to test the points routes business logic.  Unit tests can be ran by running:

```bash
mage -v unit
```

which should display all the results from the unit tests in the repository.

## Closing Thoughts

This is a jumble of a couple of different projects i've done in the past. Some things were straight copied and may look odd in the repository and others were written in a way I wouldn't describe as "produciton ready".  The end goal was to get something put together quick that could demonstrate the requirements.
