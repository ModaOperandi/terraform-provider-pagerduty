package pagerduty

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyUser() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyUserCreate,
		Read:   resourcePagerDutyUserRead,
		Update: resourcePagerDutyUserUpdate,
		Delete: resourcePagerDutyUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"email": {
				Type:     schema.TypeString,
				Required: true,
			},

			"color": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "user",
			},

			"job_title": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"avatar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"teams": {
				Type:       schema.TypeSet,
				Deprecated: "Use the 'pagerduty_team_membership' resource instead.",
				Computed:   true,
				Optional:   true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},

			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"html_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"invitation_sent": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
		},
	}
}

func buildUserStruct(d *schema.ResourceData) *pagerduty.User {
	user := &pagerduty.User{
		Name:  d.Get("name").(string),
		Email: d.Get("email").(string),
	}

	if attr, ok := d.GetOk("color"); ok {
		user.Color = attr.(string)
	}

	if attr, ok := d.GetOk("time_zone"); ok {
		user.TimeZone = attr.(string)
	}

	if attr, ok := d.GetOk("role"); ok {
		user.Role = attr.(string)
	}

	if attr, ok := d.GetOk("job_title"); ok {
		user.JobTitle = attr.(string)
	}

	if attr, ok := d.GetOk("description"); ok {
		user.Description = attr.(string)
	}

	return user
}

func resourcePagerDutyUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	user := buildUserStruct(d)

	log.Printf("[INFO] Creating PagerDuty user %s", user.Name)

	user, _, err := client.Users.Create(user)
	if err != nil {
		return err
	}

	d.SetId(user.ID)

	return resourcePagerDutyUserUpdate(d, meta)
}

func resourcePagerDutyUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Reading PagerDuty user %s", d.Id())

	user, _, err := client.Users.Get(d.Id(), &pagerduty.GetUserOptions{})
	if err != nil {
		return handleNotFoundError(err, d)
	}

	d.Set("name", user.Name)
	d.Set("email", user.Email)
	d.Set("time_zone", user.TimeZone)
	d.Set("color", user.Color)
	d.Set("role", user.Role)
	d.Set("avatar_url", user.AvatarURL)
	d.Set("description", user.Description)
	d.Set("job_title", user.JobTitle)
	d.Set("invitation_sent", user.InvitationSent)

	return nil
}

func resourcePagerDutyUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	user := buildUserStruct(d)

	log.Printf("[INFO] Updating PagerDuty user %s", d.Id())

	if _, _, err := client.Users.Update(d.Id(), user); err != nil {
		return err
	}

	return nil
}

func resourcePagerDutyUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*pagerduty.Client)

	log.Printf("[INFO] Deleting PagerDuty user %s", d.Id())

	if _, err := client.Users.Delete(d.Id()); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
