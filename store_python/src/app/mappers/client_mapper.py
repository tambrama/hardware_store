from datetime import datetime
from app.models import Client
from app.mappers.address_mapper import address_dao_to_dto


def _parse_date(v):
    if v is None:
        return None
    if hasattr(v, "isoformat"):
        return v
    if isinstance(v, str):
        return datetime.fromisoformat(v.replace("Z", "+00:00")).date()
    return v


def client_dao_to_dto(dao: Client) -> dict:
    return {
        "id": str(dao.id),
        "client_name": dao.client_name,
        "client_surname": dao.client_surname,
        "birthday": dao.birthday.isoformat() if dao.birthday else None,
        "gender": dao.gender,
        "registration_date": dao.registration_date.isoformat() if dao.registration_date else None,
        "address_id": str(dao.address_id),
        "address": address_dao_to_dto(dao.address) if dao.address else None,
    }


def client_dto_to_dao(data: dict, dao: Client = None) -> Client:
    if dao is None:
        dao = Client()
    dao.client_name = data["client_name"]
    dao.client_surname = data["client_surname"]
    dao.birthday = _parse_date(data["birthday"])
    dao.gender = data["gender"]
    dao.registration_date = _parse_date(data["registration_date"])
    dao.address_id = data["address_id"]
    return dao
