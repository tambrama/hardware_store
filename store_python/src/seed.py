from __future__ import annotations

from datetime import date
from decimal import Decimal

from app import create_app
from app.extensions import db
from app.models import Address, Client, Supplier, Product, Image


def _fake_image_bytes(name: str) -> bytes:
    # Не обязательно хранить валидный PNG/JPG: по заданию нужен byte array.
    return (f"FAKE_IMAGE::{name}").encode("utf-8")


def seed(force: bool = False) -> None:
    """
    Заполняет БД тестовыми данными.

    - По умолчанию НЕ трогает БД, если уже есть клиенты/товары/поставщики.
    - Если force=True, очищает таблицы и заполняет заново.
    """
    app = create_app()
    with app.app_context():
        db.create_all()

        if force:
            # Удаляем в правильном порядке из-за FK
            db.session.query(Product).delete()
            db.session.query(Client).delete()
            db.session.query(Supplier).delete()
            db.session.query(Image).delete()
            db.session.query(Address).delete()
            db.session.commit()
        else:
            if (
                db.session.query(Client).count() > 0
                or db.session.query(Supplier).count() > 0
                or db.session.query(Product).count() > 0
            ):
                print("Seed skipped: DB already contains data. Use --force to reseed.")
                return

        # Addresses
        addresses = [
            Address(country="Россия", city="Москва", street="ул. Ленина, 1"),
            Address(country="Россия", city="Санкт-Петербург", street="Невский пр., 10"),
            Address(country="Казахстан", city="Алматы", street="пр. Абая, 25"),
        ]
        db.session.add_all(addresses)
        db.session.commit()

        # Suppliers
        suppliers = [
            Supplier(name="ООО Поставки №1", address_id=addresses[0].id, phone_number="+7 495 111-11-11"),
            Supplier(name="ИП Дистрибьютор", address_id=addresses[1].id, phone_number="+7 812 222-22-22"),
        ]
        db.session.add_all(suppliers)
        db.session.commit()

        # Images
        images = [
            Image(image=_fake_image_bytes("fridge")),
            Image(image=_fake_image_bytes("tv")),
        ]
        db.session.add_all(images)
        db.session.commit()

        # Products
        products = [
            Product(
                name="Холодильник Arctic 3000",
                category="Холодильники",
                price=Decimal("29999.99"),
                available_stock=12,
                last_update_date=date.today(),
                supplier_id=suppliers[0].id,
                image_id=images[0].id,
            ),
            Product(
                name="Телевизор SuperVision 55\"",
                category="Телевизоры",
                price=Decimal("54990.00"),
                available_stock=7,
                last_update_date=date.today(),
                supplier_id=suppliers[1].id,
                image_id=images[1].id,
            ),
            Product(
                name="Пылесос CleanPro",
                category="Пылесосы",
                price=Decimal("8990.50"),
                available_stock=0,
                last_update_date=date.today(),
                supplier_id=suppliers[0].id,
                image_id=None,
            ),
        ]
        db.session.add_all(products)
        db.session.commit()

        # Clients
        clients = [
            Client(
                client_name="Иван",
                client_surname="Петров",
                birthday=date(1990, 5, 15),
                gender="male",
                registration_date=date.today(),
                address_id=addresses[0].id,
            ),
            Client(
                client_name="Анна",
                client_surname="Иванова",
                birthday=date(1997, 11, 3),
                gender="female",
                registration_date=date.today(),
                address_id=addresses[2].id,
            ),
        ]
        db.session.add_all(clients)
        db.session.commit()

        print("Seed done.")
        print(f"Addresses: {db.session.query(Address).count()}")
        print(f"Suppliers: {db.session.query(Supplier).count()}")
        print(f"Products: {db.session.query(Product).count()}")
        print(f"Clients: {db.session.query(Client).count()}")
        print(f"Images: {db.session.query(Image).count()}")


if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser()
    parser.add_argument("--force", action="store_true", help="Очистить таблицы и заполнить заново")
    args = parser.parse_args()

    seed(force=args.force)
