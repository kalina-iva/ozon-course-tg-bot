package messages

const manual = `Привет! Это дневник расходов.
Описание команд:
/newexpense {category} {amount} {date} - добавление нового расхода. Если дата не указана, используется текущая
/report {y|m|w} - получение отчета за последний год/месяц/неделю
`

const (
	unknownCommand        = "Неизвестная команда"
	expenseAdded          = "Расход добавлен"
	needCategoryAndAmount = "Необходимо указать категорию и сумму"
	invalidAmount         = "Некорректная сумма расхода"
	invalidDate           = "Некорректная дата"
	needPeriod            = "Необходимо указать период"
	invalidPeriod         = "Некорректный период"
	chooseCurrency        = "Выберите валюту"
	currencySaved         = "Валюта установлена"
	canNotGateRate        = "Не могу получить курс валют"
)
