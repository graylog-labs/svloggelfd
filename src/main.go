package main

import (
  "bufio"
  "io"
	"log"
	"os"
  "strings"
  "time"

	"github.com/codegangsta/cli"
  "github.com/Graylog2/go-gelf/gelf"
)

type Configuration struct {
  address  string
  source   string
  facility string
  echo     bool
  debug    bool
  tag      []string
}
var config Configuration

func sendGelf() {
  r, w := io.Pipe()
  stdoutWriter := bufio.NewWriter(os.Stdout)
  gelfWriter, err := gelf.NewWriter(config.address)
  if err != nil && config.debug {
    log.Printf("Can not create GELF connection: %s", err)
  }

  /* combine STDIN and STDERR into one io.Pipe */
  go func(writer io.Writer){
    stdinReader := bufio.NewReader(os.Stdin)
    for true {
      line, err := stdinReader.ReadString('\n')
      if err != nil {
        time.Sleep(300 * time.Millisecond)
      } else {
        writer.Write([]byte(line))
      }
    }
  }(w)
  go func(writer io.Writer){
    stderrReader := bufio.NewReader(os.Stderr)
    for true {
      line, err := stderrReader.ReadString('\n')
      if err != nil {
        time.Sleep(300 * time.Millisecond)
      } else {
        writer.Write([]byte(line))
      }
    }
  }(w)
  reader := bufio.NewReader(r)

  /* send every log line as GELF message */
  for true {
    line, err := reader.ReadString('\n')
    if err != nil {
      /* No data could be read from STDIN/STDERR */
      time.Sleep(300 * time.Millisecond)
      continue
    }

    /* additionally print log line back to STDOUT */
    if config.echo == true {
      _, err = stdoutWriter.WriteString(string(line))
      if err != nil && config.debug {
        log.Printf("Failed to write log message to STDOUT: %s", err)
      }
      err = stdoutWriter.Flush()
      if err != nil && config.debug {
        log.Printf("Failed to flush STDOUT buffer: %s", err)
      }
    }

    extraFields := make(map[string]interface{})
    if len(config.tag) == 2 {
      extraFields[getExtraFieldName(config.tag[0])] = config.tag[1]
    }

    m := gelf.Message{
      Version:  "1.1",
      Host:     config.source,
      Short:    string(line),
      TimeUnix: float64(time.Now().UnixNano()/int64(time.Millisecond)) / 1000.0,
      Level:    6,
      Facility: config.facility,
      Extra:    extraFields,
    }

    if err := gelfWriter.WriteMessage(&m); err != nil && config.debug {
      log.Printf("gelf: cannot send GELF message: %v", err)
    }
  }

  gelfWriter.Close()
}

func parsedTag(rawTag string) []string {
  return strings.Split(rawTag, ":")
}

func getExtraFieldName(key string) string {
    if key[0] == '_' {
      return key
    }
    return "_" + key
}

func main() {
  hostname, _ := os.Hostname()

	app := cli.NewApp()
	app.Name = "svgelf"
	app.Usage = "Runit GELF logger command"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
    cli.StringFlag{
			Name:      "host, H",
			Usage:     "set GELF host address and port likei: 127.0.0.1:12201",
      Value:     "localhost:12201",
		},
    cli.StringFlag{
			Name:      "source, s",
			Usage:     "override message source",
      Value:     hostname,
		},
    cli.StringFlag{
			Name:      "facility, f",
			Usage:     "override message facility",
      Value:     "runit-service",
		},
    cli.StringFlag{
			Name:      "tag, t",
			Usage:     "add a tag field to every message",
		},
    cli.BoolFlag{
			Name:      "echo, e",
			Usage:     "echo log messages back to STDOUT",
		},
    cli.BoolFlag{
      Name:      "debug, d",
      Usage:     "enable debug mode",
    },
	}
  app.Action = func(c *cli.Context) {
    config.address = c.String("host")
    config.source = c.String("source")
    config.facility = c.String("facility")
    config.tag = parsedTag(c.String("tag"))
    config.echo = c.Bool("echo")
    config.debug = c.Bool("debug")
    sendGelf()
  }
	app.Run(os.Args)
}
