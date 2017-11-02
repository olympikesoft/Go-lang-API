SELECT a.id, a.comparison_key, COALESCE(a.name, '') as name, IFNULL(a.ranking, 5+1) ranking, a.description, a.url, a.duration, a.thumb, pr.name provider, c.id category, p.price, i.label identity, t.label tag, tt.type tag_type FROM activities a JOIN (SELECT DISTINCT comparison_key, ac.ranking FROM activities ac ORDER BY ac.ranking LIMIT 5 OFFSET 0) as ck ON ck.comparison_key = a.comparison_key JOIN cities ci ON a.city_id = ci.id JOIN providers pr ON a.provider_id = pr.id JOIN pricetables p ON a.id = p.activity_id JOIN activities_categories ac ON a.id = ac.activity_id JOIN categories c ON c.id = ac.category_id JOIN activities_identities ai ON a.id = ai.activity_id JOIN identities i ON i.id = ai.identity_id JOIN activities_primary_tags apt ON a.id = apt.activity_id JOIN primary_tags pt ON pt.id = apt.primary_tag_idsJOIN activities_tags at ON a.id = at.activity_id JOIN tags t ON t.id = at.tag_id JOIN tags_types tt ON tt.id = t.type_id WHERE ci.id = 1 AND p.currency = 'EUR' GROUP BY a.id,c.id,p.price, i.name, i.label, t.label, tt.type ORDER BY a.id, a.ranking, p.price, pr.ranking, a.duration