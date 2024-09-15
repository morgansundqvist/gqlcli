# GQLCLI

GQLCLI is a command line tool for making GraphQL queries. It is designed to be portable, easy to use, and fast. It is written in Go and uses JSONATA for data manipulation / transformation.

## Why GQLCLI?

I needed a tool to make GraphQL queries and manipulate the data in a way that was easy to use and portable. I wanted to be able to run multiple queries in a single command and use data from the previous query in the next query. I also wanted to be able to filter and transform the data in a way that was easy to use and understand so I went for JSONATA as it is a powerful and easy to use query language.

I use this instead a GUI applicaiton while I'm developing and testing GraphQL queries.

## Usage

You have to create two types of files

1. GraphQL query files
2. Json configuration files

You can see examples of these files in configs/_ and graphql/_

You can run multiple queries in a single command and use data from the previous query in the next query.

The example in this repository is executed as follows:

```bash
gqlcli configs/login.json configs/company.json
```

As the variables in login.json does not exist in the context, you will be prompted to enter the values for these. This goes for missing variables in the context for all queries.

This will first login, store a company filtered by JSONATA and store the JWT token in the context. Then it will use the JWT token and companyId to query the company details.

## License

MIT
