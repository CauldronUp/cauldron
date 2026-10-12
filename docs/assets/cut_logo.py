"""Cut the Cauldron artwork off its dark background.

Solid shapes come from a BiRefNet mask; the steam and glow come from a
colour-to-alpha pass against the background, so they stay translucent.
Writes a light-text version (as drawn, for dark backgrounds) and a
dark-text version (for light backgrounds) next to this script.

    pip install numpy pillow
    python docs/assets/cut_logo.py

The mask is committed in source/. To rebuild it after changing the
artwork, delete source/cauldron-mask.png and install rembg first
(pip install "rembg[cpu]"); the script then regenerates it with the
birefnet-general model. The crop and TEXT_TOP below assume the
current 1254x1254 artwork.
"""
from pathlib import Path

import numpy as np
from PIL import Image

HERE = Path(__file__).resolve().parent
SRC = HERE / 'source' / 'cauldron-art.png'
MASK = HERE / 'source' / 'cauldron-mask.png'

if not MASK.exists():
    from rembg import new_session, remove
    remove(Image.open(SRC), session=new_session('birefnet-general'), only_mask=True).save(MASK)

c = np.asarray(Image.open(SRC).convert('RGB')).astype(np.float64)
m = np.asarray(Image.open(MASK).convert('L')).astype(np.float64) / 255

B = np.array([6.0, 10.0, 8.0])                       # sampled background
glow = np.clip((c - B) / (255 - B), 0, 1).max(axis=2)
glow = np.where(glow < 0.03, 0, glow)                # drop the vignette noise
a = np.maximum(m, glow)

safe = np.where(a > 0, a, 1)[..., None]
rgb = np.clip((c - (1 - a[..., None]) * B) / safe, 0, 255)

TEXT_TOP = 870                                       # first row below the pot legs; wordmark is below
crop = (157, 84, 1097, 1124)
size = (470, 520)

def save(rgb_arr, name):
    img = np.dstack([rgb_arr, a * 255]).round().astype(np.uint8)
    Image.fromarray(img, 'RGBA').crop(crop).resize(size, Image.LANCZOS) \
        .save(HERE / name, optimize=True)

save(rgb, 'cauldron-light.png')

dark = rgb.copy()
# recolour the wordmark only, weighted by the mask so the glow under the pot keeps its colour
w = m[TEXT_TOP:, :, None]
dark[TEXT_TOP:] = w * np.array([10.0, 10.0, 10.0]) + (1 - w) * dark[TEXT_TOP:]
save(dark, 'cauldron-dark.png')
