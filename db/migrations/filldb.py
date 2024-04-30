import psycopg2
from faker import Faker
import random
from datetime import datetime
import os
import hashlib

# Создаем экземпляр класса Faker
fake = Faker('ru_RU')

# Подключаемся к базе данных
conn = psycopg2.connect(
    dbname='socio',
    user='postgres',
    password='150600Pm',
    host='localhost'
)
cur = conn.cursor()

user_ids = []
post_ids = []
comment_ids = []

print("Connected to the database")
# Генерируем данные для таблицы user
for i in range(1000):
    # Generate a random salt
    salt = os.urandom(32)
    salt_hex = salt.hex()

    # Create a password
    password = 'password'

    # Hash the password with the salt
    key = hashlib.pbkdf2_hmac('sha256', password.encode('utf-8'), salt, 100000)
    hashed_password = key.hex()
    
    first_name = fake.first_name()
    last_name = fake.last_name()

    cur.execute(
        """
        INSERT INTO public.user (first_name, last_name, hashed_password, salt, email, avatar, date_of_birth)
        VALUES (%s, %s, %s, %s, %s, 'default_avatar.png', %s)
        RETURNING id;
        """,
        (first_name, last_name, hashed_password, salt_hex, f"email{i}@mail.ru", fake.date_of_birth())
    )
    user_ids.append(cur.fetchone()[0])

print("Data inserted into the user table")

# Генерируем данные для таблицы post
for i in range(1000000):
    if i % 10000 == 0:
        print(f"{i} posts generated")
        
    cur.execute(
        """
        INSERT INTO public.post (author_id, content, attachments)
        VALUES (%s, %s, ARRAY[%s])
        RETURNING id;
        """,
        (random.choice(user_ids), fake.text(), ', '.join(fake.words(nb=5)))
    )
    post_ids.append(cur.fetchone()[0])

print("Data inserted into the post table")

# Генерируем данные для таблицы comment
for _ in range(10000):
    cur.execute(
        """
        INSERT INTO public.comment (author_id, post_id, content)
        VALUES (%s, %s, %s)
        RETURNING id;
        """,
        (random.choice(user_ids), random.choice(post_ids), fake.text())
    )
    comment_ids.append(cur.fetchone()[0])
    
print("Data inserted into the comment table")

# Генерируем данные для таблицы personal_message
for _ in range(10000):
    cur.execute(
        """
        INSERT INTO public.personal_message (sender_id, receiver_id, content, attachments)
        VALUES (%s, %s, %s, ARRAY[%s]);
        """,
        (random.choice(user_ids), random.choice(user_ids), fake.text(), ', '.join(fake.words(nb=5)))
    )

print("Data inserted into the personal_message table")

pairs = set()
# Генерируем данные для таблицы subscription
for _ in range(1000):
    while True:
        subscriber_id = random.choice(user_ids)
        subscribed_to_id = random.choice(user_ids)
        
        # Skip if subscriber_id is the same as subscribed_to_id
        if subscriber_id == subscribed_to_id:
            continue

        # Create a pair
        pair = (subscriber_id, subscribed_to_id)

        # If the pair is not in the set, add it to the set and break the loop
        if pair not in pairs:
            pairs.add(pair)
            break

    cur.execute(
        """
        INSERT INTO public.subscription (subscriber_id, subscribed_to_id)
        VALUES (%s, %s);
        """,
        (subscriber_id, subscribed_to_id)
    )

print("Data inserted into the subscription table")

# Create a set to store the pairs
post_like_pairs = set()

for _ in range(10000):
    while True:
        post_id = random.choice(post_ids)
        user_id = random.choice(user_ids)

        # Create a pair
        pair = (post_id, user_id)

        # If the pair is not in the set, add it to the set and break the loop
        if pair not in post_like_pairs:
            post_like_pairs.add(pair)
            break

    cur.execute(
        """
        INSERT INTO post_like (post_id, user_id)
        VALUES (%s, %s);
        """,
        (post_id, user_id)
    )
    
print("Data inserted into the post_like table")

# Create a set to store the pairs
comment_like_pairs = set()

# Генерируем данные для таблицы comment_like
for _ in range(10000):
    while True:
        comment_id = random.choice(comment_ids)
        user_id = random.choice(user_ids)

        # Create a pair
        pair = (comment_id, user_id)

        # If the pair is not in the set, add it to the set and break the loop
        if pair not in comment_like_pairs:
            comment_like_pairs.add(pair)
            break

    cur.execute(
        """
        INSERT INTO comment_like (comment_id, user_id)
        VALUES (%s, %s);
        """,
        (comment_id, user_id)
    )
    
print("Data inserted into the comment_like table")

# Закрываем соединение с базой данных
conn.commit()
cur.close()
conn.close()
