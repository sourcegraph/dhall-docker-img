Small utility takes docker image strings and exports them as dhall records of parts.

Input can be a list of files, directories or stdin. Output is always stdout.

```
export DHALL_DOCKER_IMAGES=`dhall-docker-img deploy-sourcegraph/base`

dhall repl
Welcome to the Dhall v1.35.0 REPL! Type :help for more information.
⊢ :let r = env:DHALL_DOCKER_IMAGES
⊢ r

{ cadvisor =
  { name = "sourcegraph/cadvisor"
  , registry = "index.docker.io"
  , sha256 = "09076e6c5f7342de87445b295b904f28c044adb9c68b4303843fca5ddb05f832"
  , version = "insiders"
  }

  ...

, frontend =
  { name = "sourcegraph/frontend"
  , registry = "index.docker.io"
  , sha256 = "8282eed94ca7fe6b133113cef7e2ea730766abfdfaf7722a79caeed872f06ecd"
  , version = "insiders"
  }
, github-proxy =
  { name = "sourcegraph/github-proxy"
  , registry = "index.docker.io"
  , sha256 = "6222531df2a0ea88d0a8d1a3a715f0d242790b575922bededbfe2224964a893a"
  , version = "insiders"
  }
, gitserver =
  { name = "sourcegraph/gitserver"
  , registry = "index.docker.io"
  , sha256 = "a8bbb0e7ba41b812166d5df154d270801716a309fc2ff08132dcfc1c6e61d4c0"
  , version = "insiders"
  }
, syntax-highlighter =
  { name = "sourcegraph/syntax-highlighter"
  , registry = "index.docker.io"
  , sha256 = "07b9f1ff4bd2c60299f9404144cd72897fa4de2308d1be65c35bcdcd10e5410d"
  , version = "insiders"
  }
}
```
