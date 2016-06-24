## comment

Provide no rendering, i.e. ignore the contents between {% comment %} and {% endcomment %}.  The only tags that will
be parsed and acted upon are {%comment%} and {%endcomment%} - They must nest properly.

If you wish to explain why you have commented out a section you may include a set of strings after the word "comment".

Example:

```

	<p>Rendered text with {{ pub_date|date:"c" }}</p>
	{% comment "Optional Code Used created_date" "Commented out by Bob 2015-12-02" %}
		<p>Non-rendered text with {{ create_date|date:"c" }}</p>
	{% endcomment %}

```

Differences from Django: comment tags can be nested.
