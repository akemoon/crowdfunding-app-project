# GraphQL - сценарии работы

## Сценарий 1: Автор создаёт проект и следит за ним

### Шаг 1. Создать проект

```graphql
mutation {
  createProject(input: {
    userID: "01954e7a-0000-0000-0000-000000000001"
    category: "tech"
    name: "Умная теплица"
    description: "Автоматизированная система полива и мониторинга растений"
    currency: "RUB"
    goalAmount: 500000
    durationDays: 30
  })
}
```

### Шаг 2. Посмотреть свои проекты (включая те, что на модерации)

```graphql
query {
  projectsByUser(userID: "01954e7a-0000-0000-0000-000000000001") {
    id
    name
    status
    goalAmount
    currentAmount
  }
}
```

### Шаг 3. Посмотреть статус заявки на модерацию

Берём `id` проекта из шага 2.

```graphql
query {
  applicationByProject(projectID: "PROJECT_UUID") {
    status
    rejectReason
    createdAt
  }
}
```

### Шаг 4. Посмотреть проект по ID

```graphql
query {
  project(id: "PROJECT_UUID") {
    id
    name
    description
    status
    category
    goalAmount
    currentAmount
    isBoosted
    startedAt
  }
}
```

---

## Сценарий 2: Модератор обрабатывает заявку

### Шаг 1. Посмотреть все входящие заявки

```graphql
query {
  pendingApplications {
    status
    createdAt
    project {
      id
      name
      userID
      goalAmount
    }
  }
}
```

### Шаг 2. Взять заявку в работу

Берём `id` проекта из шага 1.

```graphql
mutation {
  takeApplication(
    id: "PROJECT_UUID"
    moderatorID: "01954e7a-0000-0000-0000-000000000004"
  )
}
```

### Шаг 3. Посмотреть свои заявки в работе

```graphql
query {
  myApplications(moderatorID: "01954e7a-0000-0000-0000-000000000004") {
    status
    assignedAt
    project {
      id
      name
    }
  }
}
```

### Шаг 4а. Одобрить проект

```graphql
mutation {
  approveProject(id: "PROJECT_UUID")
}
```

### Шаг 4б. Или отклонить с причиной

```graphql
mutation {
  rejectProject(
    id: "PROJECT_UUID"
    reason: "Недостаточно информации о команде проекта"
  )
}
```

---

## Сценарий 3: Просмотр активных проектов и буст

### Шаг 1. Список активных проектов

```graphql
query {
  projects(filter: { status: "active", limit: 10 }) {
    id
    name
    goalAmount
    currentAmount
    isBoosted
    startedAt
  }
}
```

### Шаг 2. Фильтр по категории и поиск

```graphql
query {
  projects(filter: {
    status: "active"
    category: "tech"
    search: "теплица"
  }) {
    id
    name
    currentAmount
    goalAmount
  }
}
```

### Шаг 3. Забустить проект промокодом

```graphql
mutation {
  boostProject(
    id: "PROJECT_UUID"
    userID: "01954e7a-0000-0000-0000-000000000001"
    promoCode: "BOOST2024"
  )
}
```

### Шаг 4. Убедиться, что проект стал boosted

```graphql
query {
  project(id: "PROJECT_UUID") {
    name
    isBoosted
    boostedUntil
  }
}
```

---

## Фишка GraphQL: несколько запросов за один раз

```graphql
query {
  pending: pendingApplications {
    status
    project { id name }
  }

  activeProjects: projects(filter: { status: "active", limit: 5 }) {
    id
    name
    isBoosted
  }
}
```
