package config

var (
	PROMPT_SYSTEM    = "You are a program that classifies data from OCR receipts"
	PROMPT_ASSISTANT = "Classify the data of the receipts you are about to post in Japanese. Types: 1=food, 2=drink, 3=ingredients, 9=alcohol/tobacco, 0=other. No extra explanation.\n"
	PROMPT_USER      = "Convert the above receipt data into the same format as the sample below. No extra explanation.\n\n"
)

var PROMPT_EXAMPLE string = `{
	"store": {"name": "Name", "address": "Address", "telephone": "Phone"},
	"datetime": "2023-06-10T12:34:56",
	"items": [
		{"item_name": "Item1", "quantity": 1, "unit_price": 100, "type": 1},
		{"item_name": "Item2", "quantity": 2, "unit_price": 200, "type": 0}
	],
	"tax": 50,
	"tax_included_price": 550
}`
