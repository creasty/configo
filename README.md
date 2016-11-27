configo
=======

Definitive configuration utility in Golang


What
----

`configo` gathers configuration parameters from variety of sources, overlays values by a different env -- production, development, etc -- and merges into a single struct.

### Features

It's a set of popular libraries.

- Read from environment vars / dotenv file and bind them into struct
  - https://github.com/joho/godotenv
  - https://github.com/kelseyhightower/envconfig
- Read from YAML files
  - https://github.com/ghodss/yaml
- Finally, validate the struct
  - https://github.com/asaskevich/govalidator

### Examples

Env vars:

```sh
export APP_ENV=production
```

Dotenv file:

```sh
# .env

APP_FOO=foo
```

Config yaml files:

```yaml
# config/default.yml

Bar: default
Baz: default?
```

```yaml
# config/production.yml

Baz: production
```

And your config struct...

```go
type Config struct {
	Env string
	Foo string
	Bar string
	Baz string
}

c := Config{}
configo.Load(c, configo.Option{Dir: "./config"})

c.Env //=> "production"
c.Foo //=> "foo"
c.Bar //=> "default"
c.Baz //=> "production"
```


License
-------

This project is copyright by [Creasty](http://creasty.com), released under the MIT license.  
See `LICENSE.txt` file for details.
