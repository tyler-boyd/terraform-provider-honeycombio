package honeycombio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/kvrhdn/go-honeycombio"
)

func newQueryAnnotation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceQueryAnnotationCreate,
		ReadContext:   resourceQueryAnnotationRead,
		UpdateContext: resourceQueryAnnotationUpdate,
		DeleteContext: schema.NoopContext,
		Importer:      nil,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// "query_id": {
			// 	Type:     schema.TypeString,
			// 	Required: true,
			// 	ForceNew: true,
			// },
			"dataset": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceQueryAnnotationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*honeycombio.Client)

	dataset := d.Get("dataset").(string)
	queryAnnotation := readQueryAnnotation(d)

	query, err := client.Queries.Create(ctx, dataset, &honeycombio.QuerySpec{})

	if err != nil {
		return diag.FromErr(err)
	}

	queryAnnotation.QueryID = *query.ID

	queryAnnotation, err = client.QueryAnnotations.Create(ctx, dataset, queryAnnotation)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(queryAnnotation.ID)

	return resourceQueryAnnotationRead(ctx, d, meta)
}

func resourceQueryAnnotationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*honeycombio.Client)

	dataset := d.Get("dataset").(string)

	queryAnnotation, err := client.QueryAnnotations.Get(ctx, dataset, d.Id())
	if err != nil {
		if err == honeycombio.ErrNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.SetId(queryAnnotation.ID)
	d.Set("name", queryAnnotation.Name)
	d.Set("description", queryAnnotation.Description)
	d.Set("query_id", queryAnnotation.QueryID)
	return nil
}

func resourceQueryAnnotationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*honeycombio.Client)

	dataset := d.Get("dataset").(string)
	queryAnnotation := readQueryAnnotation(d)

	_, err := client.QueryAnnotations.Update(ctx, dataset, queryAnnotation)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceQueryAnnotationRead(ctx, d, meta)
}

func readQueryAnnotation(d *schema.ResourceData) *honeycombio.QueryAnnotation {
	return &honeycombio.QueryAnnotation{
		ID:          d.Id(),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}
}
