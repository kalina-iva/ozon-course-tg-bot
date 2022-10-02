package messages

const manual = `Привет! Это дневник расходов.
Описание команд:
/newcat {name} - добавление новой категории
/allcat - просмотр всех категории
/newexpense {categoryNumber} {amount} {date} - добавление нового расхода. Если дата не указана, используется текущая
/report {y|m|w} - получение отчета за последний год/месяц/неделю
`

const unknownCommand = "Неизвестная команда"
const noCategories = "Нет категорий"
const expenseAdded = "Расход добавлен"
const needCategoryAndAmount = "Необходимо указать категорию и сумму"
const invalidCategoryNumber = "Некорректный номер категории"
const invalidAmount = "Некорректная сумма расхода"
const invalidDate = "Некорректная дата"
const needPeriod = "Необходимо указать период"
const invalidPeriod = "Некорректный период"
