package views

var Views = map[string]string{
	"vw_organization_member_roles": `create view vw_organization_member_roles as
select o.id as organization_id, o.name as organization_name, o.slug, om.user_id, omr.* from organization_members om 
left join organizations o on om.organization_id = o.id 
left join organization_member_roles omr on om.role_id = omr.id `,
}
