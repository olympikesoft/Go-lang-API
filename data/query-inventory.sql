SELECT
a.id,
pr.name provider,
c.name category_name,
c.id category_id,
c.ui_priority category_priority,
pt.name tag,
i.name identity,
i.type identity_type,
oper.name operators,
p.price
FROM activities a
JOIN operators oper ON a.operator_id = oper.id
JOIN cities ci ON a.city_id = ci.id
JOIN providers pr ON a.provider_id = pr.id
JOIN pricetables p ON a.id = p.activity_id
JOIN activities_categories ac ON a.id = ac.activity_id
JOIN categories c ON c.id = ac.category_id
JOIN activities_identities ai ON a.id = ai.activity_id
JOIN identities i ON i.id = ai.identity_id
JOIN activities_primary_tags apt ON a.id = apt.activity_id
JOIN primary_tags pt ON pt.id = apt.primary_tag_id
WHERE ci.id =1
AND p.currency = 'EUR'
ORDER by p.price, a.duration