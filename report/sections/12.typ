= Творческая задача: игра "DODEP RP"
== Идея

Генерируется случайный связный граф. Игроки спавнятся в разных углах и должны собирать монетки на протяжении 45 секунд. Победителем считается тот, кто

== Код

Описание класса Coin

```py
class Coin:
    def __init__(self, node_index, nodes):
        self.node_index = node_index
        self.x, self.y = nodes[node_index]
        self.collected = False
        self.spawn_time = 0

    def draw(self, screen, camera_offset=(0, 0)):
        if not self.collected:
            import pygame

            x = self.x + camera_offset[0]
            y = self.y + camera_offset[1]
            pygame.draw.circle(screen, (255, 215, 0), (int(x), int(y)), 8)
            pygame.draw.circle(screen, (255, 255, 0), (int(x), int(y)), 8, 2)

```

Описание класса GraphGenerator

```py
import math
import random


class GraphGenerator:
    def __init__(self, width, height, node_count=25):
        self.width = width
        self.height = height
        self.node_count = node_count
        self.nodes = []
        self.edges = []

    def generate_grid_for_directions(self):
        """Генерация графа, удобного для управления по направлениям"""
        print("Генерируем граф для управления стрелками...")

        self.nodes = []
        self.edges = []

        # Создаём сетку 6x8
        rows = 6
        cols = 8
        spacing_x = 180
        spacing_y = 150

        for row in range(rows):
            for col in range(cols):
                # Центрируем узел в ячейке
                x = 150 + col * spacing_x
                y = 150 + row * spacing_y

                # Небольшое случайное смещение для естественности
                x += random.randint(-20, 20)
                y += random.randint(-20, 20)

                self.nodes.append((x, y))

        print(f"Создано {len(self.nodes)} узлов в сетке")

        # Создаём основные рёбра для 4-х направлений
        for i in range(len(self.nodes)):
            row_i = i // cols
            col_i = i % cols

            # До 4 соседей: вверх, вниз, влево, вправо
            neighbors = []

            # Вправо
            if col_i < cols - 1:
                neighbors.append(i + 1)

            # Вниз
            if row_i < rows - 1:
                neighbors.append(i + cols)

            # Влево
            if col_i > 0:
                neighbors.append(i - 1)

            # Вверх
            if row_i > 0:
                neighbors.append(i - cols)

            # Ограничиваем максимум 4 соседей
            if len(neighbors) > 4:
                # Выбираем случайных 4 соседей
                neighbors = random.sample(neighbors, 4)

            # Добавляем рёбра
            for neighbor in neighbors:
                if i < neighbor:  # Добавляем ребро только один раз
                    self.edges.append((i, neighbor))

        # Удаляем некоторые рёбра для разнообразия (но оставляем связность)
        edges_to_remove = []
        for u, v in self.edges:
            # Проверяем степень вершин
            degree_u = sum(1 for x, y in self.edges if x == u or y == u)
            degree_v = sum(1 for x, y in self.edges if x == v or y == v)

            # Удаляем некоторые рёбра с вероятностью
            if random.random() < 0.2 and degree_u > 2 and degree_v > 2:
                edges_to_remove.append((u, v))

        for edge in edges_to_remove:
            if edge in self.edges:
                self.edges.remove(edge)

        # Проверяем связность
        if not self.is_connected():
            print("Восстанавливаем связность...")
            self.ensure_connectivity()

        print(f"Создано {len(self.edges)} рёбер")
        self.print_statistics()

        return self.nodes, self.edges

    def generate(self):
        """Основной метод генерации"""
        return self.generate_grid_for_directions()

    def get_neighbors(self, node_index):
        """Получить всех соседей узла"""
        neighbors = []
        for u, v in self.edges:
            if u == node_index:
                neighbors.append(v)
            elif v == node_index:
                neighbors.append(u)
        return neighbors

    def is_connected(self):
        """Проверка связности графа"""
        if not self.nodes:
            return False

        visited = [False] * len(self.nodes)
        stack = [0]

        while stack:
            node = stack.pop()
            if not visited[node]:
                visited[node] = True
                neighbors = self.get_neighbors(node)
                stack.extend(neighbors)

        return all(visited)

    def ensure_connectivity(self):
        """Гарантируем связность графа"""
        visited = [False] * len(self.nodes)
        components = []

        for i in range(len(self.nodes)):
            if not visited[i]:
                component = []
                stack = [i]

                while stack:
                    node = stack.pop()
                    if not visited[node]:
                        visited[node] = True
                        component.append(node)
                        neighbors = self.get_neighbors(node)
                        stack.extend(neighbors)

                components.append(component)

        # Соединяем компоненты
        for i in range(len(components) - 1):
            # Берём случайные узлы из соседних компонент
            node1 = random.choice(components[i])
            node2 = random.choice(components[i + 1])
            self.edges.append((node1, node2))

    def print_statistics(self):
        """Вывод статистики графа"""
        print("\n" + "=" * 50)
        print("СТАТИСТИКА ГРАФА:")
        print("=" * 50)

        degree_count = {}
        max_degree = 0
        min_degree = float("inf")

        for i in range(len(self.nodes)):
            degree = len(self.get_neighbors(i))
            degree_count[degree] = degree_count.get(degree, 0) + 1
            max_degree = max(max_degree, degree)
            min_degree = min(min_degree, degree)

        print(f"Узлов: {len(self.nodes)}")
        print(f"Рёбер: {len(self.edges)}")
        print(f"Минимальная степень: {min_degree}")
        print(f"Максимальная степень: {max_degree}")
        print(f"Распределение степеней: {dict(sorted(degree_count.items()))}")

        # Проверяем ограничения
        violations = []
        for i in range(len(self.nodes)):
            if len(self.get_neighbors(i)) > 4:
                violations.append(f"Узел {i}: степень {len(self.get_neighbors(i))} > 4")

        if violations:
            print("НАРУШЕНИЯ ОГРАНИЧЕНИЙ:")
            for violation in violations:
                print(f"  - {violation}")
        else:
            print("✓ Все ограничения соблюдены!")
        print("=" * 50)

```

Описание класса Player

```py
import math

import pygame


class Player:
    def __init__(self, start_node, color, controls, name):
        self.current_node = start_node
        self.color = color
        self.score = 0
        self.controls = controls  # {'up': pygame.K_UP, 'down': pygame.K_DOWN, 'left': pygame.K_LEFT, 'right': pygame.K_RIGHT}
        self.name = name

        # Для плавного движения
        self.x, self.y = 0, 0
        self.target_x, self.target_y = 0, 0
        self.is_moving = False
        self.move_progress = 0
        self.move_speed = 0.25  # Быстрое движение
        self.move_from_node = start_node
        self.move_to_node = start_node

        # Для управления
        self.last_move_time = 0
        self.move_cooldown = 100

        # Для отображения сообщений
        self.blocked_message = ""
        self.blocked_message_time = 0

    def find_neighbor_in_direction(
        self, neighbors, nodes, direction, opponent_node=None
    ):
        """Найти соседа в заданном направлении (up, down, left, right)"""
        if not neighbors:
            return None

        current_x, current_y = self.x, self.y
        best_neighbor = None
        best_angle_diff = float("inf")

        direction_angles = {
            "right": 0,
            "up": 90,
            "left": 180,
            "down": 270,
        }

        target_angle = direction_angles[direction]

        for neighbor in neighbors:
            # Проверяем, не занят ли этот узел соперником
            if opponent_node is not None and neighbor == opponent_node:
                continue

            nx, ny = nodes[neighbor]
            dx = nx - current_x
            dy = ny - current_y

            angle = math.degrees(
                math.atan2(-dy, dx)
            )
            if angle < 0:
                angle += 360

            angle_diff = min(abs(angle - target_angle), 360 - abs(angle - target_angle))

            if angle_diff < best_angle_diff:
                best_angle_diff = angle_diff
                best_neighbor = neighbor

        if best_angle_diff <= 60:
            return best_neighbor
        return None

    def move_in_direction(self, direction, nodes, neighbors, opponent_node=None):
        """Начать движение в заданном направлении"""
        if not self.is_moving:
            neighbor = self.find_neighbor_in_direction(
                neighbors, nodes, direction, opponent_node
            )

            if neighbor is not None:
                self.move_from_node = self.current_node
                self.move_to_node = neighbor
                self.x, self.y = nodes[self.current_node]
                self.target_x, self.target_y = nodes[neighbor]
                self.is_moving = True
                self.move_progress = 0
                self.last_move_time = pygame.time.get_ticks()
                self.blocked_message = ""
                return True
            else:
                if opponent_node is not None:
                    for n in neighbors:
                        if n == opponent_node:
                            nx, ny = nodes[n]
                            dx = nx - self.x
                            dy = ny - self.y
                            angle = math.degrees(math.atan2(-dy, dx))
                            if angle < 0:
                                angle += 360

                            direction_angles = {
                                "right": 0,
                                "up": 90,
                                "left": 180,
                                "down": 270,
                            }
                            target_angle = direction_angles[direction]
                            angle_diff = min(
                                abs(angle - target_angle),
                                360 - abs(angle - target_angle),
                            )

                            if angle_diff <= 60:
                                self.blocked_message = "BLOCKED!"
                                self.blocked_message_time = pygame.time.get_ticks()
                                print(
                                    f"{self.name}: направление {direction} заблокировано соперником!"
                                )
                                return False

                return False
        return False

    def update(self, nodes, current_time):
        """Обновление положения"""
        if self.is_moving:
            self.move_progress += self.move_speed
            if self.move_progress >= 1:
                self.move_progress = 1
                self.is_moving = False
                self.current_node = self.move_to_node
                self.x, self.y = self.target_x, self.target_y
            else:
                # Линейная интерполяция
                t = self.move_progress
                self.x = self.x + (self.target_x - self.x) * t
                self.y = self.y + (self.target_y - self.y) * t

        # Очищаем сообщение о блокировке через 1 секунду
        if self.blocked_message and current_time - self.blocked_message_time > 1000:
            self.blocked_message = ""

    def get_position(self):
        """Получить текущую позицию"""
        return self.x, self.y

    def draw(self, screen, nodes, neighbors, camera_offset=(0, 0), font=None):
        """Отрисовка игрока"""
        x, y = self.x, self.y

        # Рисуем направляющие линии к соседям
        if not self.is_moving and neighbors:
            for neighbor in neighbors:
                nx, ny = nodes[neighbor]

                dx = nx - x
                dy = ny - y
                dist = math.sqrt(dx * dx + dy * dy)
                if dist > 0:
                    dx, dy = dx / dist, dy / dist

                    # Определяем направление
                    angle = math.degrees(math.atan2(-dy, dx))
                    if angle < 0:
                        angle += 360

                    # Цвет линии в зависимости от направления
                    if 45 <= angle < 135:  # Вверх
                        line_color = (100, 255, 100)  # Зелёный
                    elif 135 <= angle < 225:  # Влево
                        line_color = (100, 100, 255)  # Синий
                    elif 225 <= angle < 315:  # Вниз
                        line_color = (255, 100, 100)  # Красный
                    else:  # Вправо
                        line_color = (255, 255, 100)  # Жёлтый

                    # Рисуем линию к соседу
                    arrow_x = x + dx * 40
                    arrow_y = y + dy * 40

                    pygame.draw.line(
                        screen,
                        line_color,
                        (int(x + dx * 25), int(y + dy * 25)),
                        (int(arrow_x), int(arrow_y)),
                        2,
                    )

                    # Рисуем стрелку
                    arrow_angle = math.atan2(dy, dx)
                    arrow_size = 8

                    points = []
                    for i in range(3):
                        angle_i = arrow_angle + math.pi * 0.8 + i * math.pi * 0.8
                        px = arrow_x + math.cos(angle_i) * arrow_size
                        py = arrow_y + math.sin(angle_i) * arrow_size
                        points.append((px, py))

                    pygame.draw.polygon(screen, line_color, points)

        # Рисуем игрока
        pygame.draw.circle(screen, self.color, (int(x), int(y)), 22)
        pygame.draw.circle(screen, (255, 255, 255), (int(x), int(y)), 22, 3)

        # Внутренний круг
        inner_color = (
            min(self.color[0] + 70, 255),
            min(self.color[1] + 70, 255),
            min(self.color[2] + 70, 255),
        )
        pygame.draw.circle(screen, inner_color, (int(x), int(y)), 10)

        # Рисуем значок "занято" если игрок стоит на узле
        pygame.draw.circle(screen, (255, 255, 255, 100), (int(x), int(y)), 28, 2)

        # Отображение имени и счёта
        if font:
            # Имя над игроком
            name_text = font.render(self.name, True, (255, 255, 255))
            name_rect = name_text.get_rect(center=(int(x), int(y - 40)))
            screen.blit(name_text, name_rect)

            # Счёт под игроком
            score_text = font.render(str(self.score), True, (255, 255, 255))
            score_rect = score_text.get_rect(center=(int(x), int(y + 40)))
            screen.blit(score_text, score_rect)

            # Сообщение о блокировке
            if self.blocked_message:
                blocked_text = font.render(self.blocked_message, True, (255, 50, 50))
                blocked_rect = blocked_text.get_rect(center=(int(x), int(y - 70)))
                screen.blit(blocked_text, blocked_rect)

```

Содержимое файла `ui.py`

```py
import pygame


class Button:
    def __init__(
        self,
        x,
        y,
        width,
        height,
        text,
        color=(100, 100, 200),
        hover_color=(120, 120, 220),
        text_color=(255, 255, 255),
    ):
        self.rect = pygame.Rect(x, y, width, height)
        self.text = text
        self.color = color
        self.hover_color = hover_color
        self.text_color = text_color
        self.is_hovered = False
        self.font = pygame.font.Font(None, 36)

    def draw(self, screen):
        color = self.hover_color if self.is_hovered else self.color
        pygame.draw.rect(screen, color, self.rect, border_radius=12)
        pygame.draw.rect(screen, (200, 200, 255), self.rect, 3, border_radius=12)

        text_surf = self.font.render(self.text, True, self.text_color)
        text_rect = text_surf.get_rect(center=self.rect.center)
        screen.blit(text_surf, text_rect)

    def check_hover(self, pos):
        self.is_hovered = self.rect.collidepoint(pos)
        return self.is_hovered

    def is_clicked(self, pos, event):
        if event.type == pygame.MOUSEBUTTONDOWN and event.button == 1:
            return self.rect.collidepoint(pos)
        return False


class InputBox:
    def __init__(self, x, y, width, height, text="", max_length=10):
        self.rect = pygame.Rect(x, y, width, height)
        self.text = text
        self.max_length = max_length
        self.active = False
        self.font = pygame.font.Font(None, 32)
        self.color_inactive = (80, 80, 120)
        self.color_active = (100, 100, 180)
        self.color = self.color_inactive

    def handle_event(self, event):
        if event.type == pygame.MOUSEBUTTONDOWN:
            if self.rect.collidepoint(event.pos):
                self.active = not self.active
            else:
                self.active = False
            self.color = self.color_active if self.active else self.color_inactive

        if event.type == pygame.KEYDOWN:
            if self.active:
                if event.key == pygame.K_RETURN:
                    return True
                elif event.key == pygame.K_BACKSPACE:
                    self.text = self.text[:-1]
                else:
                    if len(self.text) < self.max_length:
                        self.text += event.unicode
        return False

    def update(self):
        pass

    def draw(self, screen):
        pygame.draw.rect(screen, self.color, self.rect, border_radius=8)
        pygame.draw.rect(screen, (200, 200, 255), self.rect, 2, border_radius=8)

        text_surf = self.font.render(self.text, True, (255, 255, 255))
        screen.blit(text_surf, (self.rect.x + 10, self.rect.y + 10))

        if self.active:
            cursor_pos = self.font.size(self.text)[0] + self.rect.x + 12
            pygame.draw.line(
                screen,
                (255, 255, 255),
                (cursor_pos, self.rect.y + 8),
                (cursor_pos, self.rect.y + self.rect.height - 8),
                2,
            )


class Menu:
    def __init__(self, title, width, height):
        self.width = width
        self.height = height
        self.title = title
        self.buttons = []
        self.input_boxes = []
        self.title_font = pygame.font.Font(None, 72)
        self.font = pygame.font.Font(None, 36)

    def draw(self, screen):
        # Полупрозрачный фон
        overlay = pygame.Surface((self.width, self.height), pygame.SRCALPHA)
        overlay.fill((0, 0, 0, 200))
        screen.blit(overlay, (0, 0))

        # Заголовок
        title_text = self.title_font.render(self.title, True, (255, 255, 255))
        title_rect = title_text.get_rect(center=(self.width // 2, self.height // 4))
        screen.blit(title_text, title_rect)

        # Кнопки
        for button in self.buttons:
            button.draw(screen)

        # Поля ввода
        for input_box in self.input_boxes:
            input_box.draw(screen)

```

== Запуск игры

#image("/assets/image-1.png")

#image("/assets/image-2.png")

#image("/assets/image-3.png")

#image("/assets/image-4.png")

#image("/assets/image-5.png")

== Описание алгоритма

Алгоритм генерации графа состоит из трёх этапов.

*Этап 1. Создание узлов.*
На поле размещается упорядоченная сетка узлов размером 6 на 8. Каждому узлу назначаются координаты с равным интервалом. Для визуального разнообразия координаты слегка смещаются в случайном направлении.

*Этап 2. Создание рёбер.*
Каждый узел соединяется рёбрами с соседними узлами справа, слева, сверху и снизу. Это гарантирует базовую связность и обеспечивает возможность перемещения в четырёх направлениях.

*Этап 3. Оптимизация структуры.*
Часть рёбер удаляется для создания более разнообразных путей, но степень каждого узла сохраняется не менее двух. Проверяется связность графа: если образуются изолированные группы узлов, они соединяются дополнительными рёбрами.

Результатом является связный граф, где из каждого узла можно переместиться в соседний по одному из четырёх направлений, что удобно для управления стрелками клавиатуры.