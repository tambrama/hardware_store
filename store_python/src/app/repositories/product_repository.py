import uuid
from app.extensions import db
from app.models import Product


class ProductRepository:
    """Репозиторий для работы с товарами (DAO)."""

    def get_by_id(self, product_id: uuid.UUID) -> Product | None:
        return db.session.get(Product, product_id)

    def get_all_available(self) -> list[Product]:
        return list(
            db.session.query(Product)
            .filter(Product.available_stock > 0)
            .order_by(Product.name)
            .all()
        )

    def add(self, product: Product) -> Product:
        db.session.add(product)
        db.session.commit()
        return product

    def update(self, product: Product) -> Product:
        db.session.commit()
        return product

    def delete(self, product: Product) -> None:
        db.session.delete(product)
        db.session.commit()
