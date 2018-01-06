import os

import magic
from django.db import models
from django.db.models.signals import post_delete, pre_save
from django.dispatch import receiver


class Object(models.Model):
    """
    스토리지 파일
    """

    file = models.FileField()
    content_type = models.CharField(max_length=200, default='text')

    name = models.CharField(max_length=200, unique=True)
    uploader = models.CharField(max_length=200)
    safety_checked = models.BooleanField(default=False)

    created_at = models.DateTimeField(auto_now_add=True)
    modified_at = models.DateTimeField(auto_now=True)


@receiver(pre_save, sender=Object)
def pre_object_save(sender, instance, *args, **kwargs):
    instance.content_type = magic.from_buffer(instance.file.read(1024), mime=True)


@receiver(post_delete, sender=Object)
def post_object_delete(sender, instance, *args, **kwargs):
    if os.path.isfile(instance.file.path):
        os.remove(instance.file.path)
