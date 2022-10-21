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
	canNotGetRate         = "Не могу получить курс валют"
	canNotSaveCurrency    = "Произошла ошибка при установлении валюты"
	limitSaved            = "Лимит установлен"
	canNotSaveLimit       = "Возникла ошибка при изменении лимита"
	limitDeleted          = "Лимит сброшен"
	canNotAddExpense      = "Произошла ошибка при добавлении расхода"
	canNotCreateReport    = "Произошла ошибка при формировании отчета"
	limitExceeded         = "Превышен лимит"
)
