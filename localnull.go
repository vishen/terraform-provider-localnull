package main

import (
	"bytes"
	"html/template"
	"log"
	"os/exec"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLocalNull() *schema.Resource {
	return &schema.Resource{
		Create: resourceServerCreate,
		Read:   resourceServerRead,
		Update: resourceServerUpdate,
		Delete: resourceServerDelete,

		Schema: map[string]*schema.Schema{
			"shell": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "bash",
			},
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Output variables
			"interpolated_command": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"output": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerCreate(d *schema.ResourceData, m interface{}) error {
	return handleCommand(d)
}

func handleCommand(d *schema.ResourceData) error {
	command, err := interpolatedCommand(d)
	if err != nil {
		return err
	}

	shell := d.Get("shell").(string)

	cmd := exec.Command(shell, "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	d.SetId(shell + command)
	if err := d.Set("interpolated_command", command); err != nil {
		return err
	}
	if err := d.Set("output", out.String()); err != nil {
		return err
	}
	return nil
}

func interpolatedCommand(d *schema.ResourceData) (string, error) {
	configuration := d.Get("configuration").(map[string]interface{})
	command := d.Get("command").(string)

	tmpl, err := template.New("command-with-variables").Parse(command)
	if err != nil {
		return "", err
	}
	buf := bytes.Buffer{}
	if err := tmpl.Execute(&buf, configuration); err != nil {
		return "", err
	}
	log.Printf("[DEBUG] GET %q + %v == %q\n", command, configuration, buf.String())
	return buf.String(), nil
}

func resourceServerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServerUpdate(d *schema.ResourceData, m interface{}) error {
	return handleCommand(d)
}

func resourceServerDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
