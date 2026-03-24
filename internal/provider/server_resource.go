package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = (*serverResource)(nil)
var _ resource.ResourceWithConfigure = (*serverResource)(nil)
var _ resource.ResourceWithImportState = (*serverResource)(nil)

func NewServerResource() resource.Resource {
	return &serverResource{}
}

type serverResource struct {
	client *Client
}

type serverResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Pool        types.String `tfsdk:"pool"`
	Hostname    types.String `tfsdk:"hostname"`
	Description types.String `tfsdk:"description"`
	Cores       types.Int64  `tfsdk:"cores"`
	Memory      types.Int64  `tfsdk:"memory"`
	Disks       types.List   `tfsdk:"disks"`
	Image       types.String `tfsdk:"image"`
	CustomName  types.String `tfsdk:"custom_name"`
	Password    types.String `tfsdk:"password"`
	SSHKey      types.String `tfsdk:"sshkey"`
	Type        types.String `tfsdk:"type"`
	Subnets     types.List   `tfsdk:"subnets"`
	SubnetsInt  types.List   `tfsdk:"subnets_intern"`
	Tags        types.List   `tfsdk:"tags"`
	Rollouts    types.List   `tfsdk:"rollouts"`
}

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Centron Cloud Server.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the server.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pool": schema.StringAttribute{
				Description: "Pool name.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname of the server.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the server.",
				Optional:    true,
			},
			"cores": schema.Int64Attribute{
				Description: "Number of CPU cores.",
				Required:    true,
			},
			"memory": schema.Int64Attribute{
				Description: "Memory in MB.",
				Required:    true,
			},
			"disks": schema.ListAttribute{
				Description: "List of Disk capacities.",
				ElementType: types.Int64Type,
				Required:    true,
			},
			"image": schema.StringAttribute{
				Description: "OS Image string identifier.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"custom_name": schema.StringAttribute{
				Description: "Custom Name for the server.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Server access password.",
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"sshkey": schema.StringAttribute{
				Description: "SSH key to inject.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				Description: "Either generic, managed, unmanaged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subnets": schema.ListAttribute{
				Description: "List of public subnet IDs.",
				ElementType: types.Int64Type,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"subnets_intern": schema.ListAttribute{
				Description: "List of internal subnet IDs.",
				ElementType: types.Int64Type,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"tags": schema.ListAttribute{
				Description: "List of tags.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"rollouts": schema.ListAttribute{
				Description: "Rollouts applied on creation.",
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *serverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan serverResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := map[string]interface{}{
		"pool":      plan.Pool.ValueString(),
		"hostname":  plan.Hostname.ValueString(),
		"cores":     plan.Cores.ValueInt64(),
		"memory":    plan.Memory.ValueInt64(),
		"image":     plan.Image.ValueString(),
	}

	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		payload["description"] = plan.Description.ValueString()
	}
	if !plan.CustomName.IsNull() && !plan.CustomName.IsUnknown() {
		payload["custom_name"] = plan.CustomName.ValueString()
	}
	if !plan.Password.IsNull() && !plan.Password.IsUnknown() {
		payload["password"] = plan.Password.ValueString()
	}
	if !plan.SSHKey.IsNull() && !plan.SSHKey.IsUnknown() {
		payload["sshkey"] = plan.SSHKey.ValueString()
	}
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		payload["type"] = plan.Type.ValueString()
	}

	// Disks
	if !plan.Disks.IsNull() && !plan.Disks.IsUnknown() {
		var disks []int64
		diags := plan.Disks.ElementsAs(ctx, &disks, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		payload["disks"] = disks
	}

	// Subnets
	if !plan.Subnets.IsNull() && !plan.Subnets.IsUnknown() {
		var subnets []int64
		diags := plan.Subnets.ElementsAs(ctx, &subnets, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(subnets) > 0 {
			payload["subnets"] = subnets
		}
	}

	// Subnets Intern
	if !plan.SubnetsInt.IsNull() && !plan.SubnetsInt.IsUnknown() {
		var subnetsInt []int64
		diags := plan.SubnetsInt.ElementsAs(ctx, &subnetsInt, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(subnetsInt) > 0 {
			payload["subnets_intern"] = subnetsInt
		}
	}

	// Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		diags := plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(tags) > 0 {
			payload["tags"] = tags
		}
	}

	// Rollouts
	if !plan.Rollouts.IsNull() && !plan.Rollouts.IsUnknown() {
		var rollouts []string
		diags := plan.Rollouts.ElementsAs(ctx, &rollouts, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		if len(rollouts) > 0 {
			payload["rollouts"] = rollouts
		}
	}

	type CreateResp struct {
		ID int `json:"id"`
	}

	var createResp CreateResp
	err := r.client.DoRequest(ctx, "POST", "/ccloud/servers", payload, &createResp)
	if err != nil {
		resp.Diagnostics.AddError("Error creating server", err.Error())
		return
	}
	// Polling for server state = provisioned to be added in future iterations.
	plan.ID = types.StringValue(strconv.Itoa(createResp.ID))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state serverResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	type ServerResp struct {
		Data struct {
			Cores  int    `json:"cores"`
			Memory int    `json:"memory"`
			Type   string `json:"type"`
			// future expansion fields here
		} `json:"data"`
	}

	var sr ServerResp
	// GET /ccloud/servers/{hostname}
	err := r.client.DoRequest(ctx, "GET", "/ccloud/servers/"+state.Hostname.ValueString(), nil, &sr)
	if err != nil {
		resp.Diagnostics.AddError("Error reading server", err.Error())
		return
	}

	// Sync State with retrieved values. Note: GET response payload differs from POST.
	state.Cores = types.Int64Value(int64(sr.Data.Cores))
	state.Memory = types.Int64Value(int64(sr.Data.Memory))
	if sr.Data.Type != "" {
		state.Type = types.StringValue(sr.Data.Type)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Complex updates, like scaling memory/cores without recreation, are reserved for future iterations.
	var plan serverResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state serverResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// DELETE /ccloud/servers/{hostname}
	err := r.client.DoRequest(ctx, "DELETE", "/ccloud/servers/"+state.Hostname.ValueString(), nil, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting server", err.Error())
		return
	}
}

func (r *serverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("hostname"), req, resp)
}
