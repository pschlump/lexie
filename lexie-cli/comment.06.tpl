
Before Comment {{ DictionarySubstituteBefore }}
{% comment "Optional note" %}
    <p>Commented out text with {{ create_date|date:"c" }}</p>
{% endcomment %}
After Comment {{ DictionarySubstituteAfter }}
-- Documentation notes that comments can not be nested ? why?
https://docs.djangoproject.com/en/1.8/ref/templates/builtins/#ref-templates-builtins-tags
