package openai

import "strings"

const prompt = `Redesign this interior photo into a refined and super creative [STYLE] [ROOM_TYPE].
Do not change architecture:
Keep layout, geometry, perspective, proportions, windows, doors, and walls exactly as in the original.
No new or moved structural elements.
Enhance only interiors:
Remove clutter, construction items, and old furniture.
Add elegant, imaginative furniture, lighting, textures, and décor in the chosen style.
Be bold, artistic, and visually striking — but stay harmonious.
Result must look like a polished, magazine-quality photoshoot with natural lighting and realistic materials.
Output: photorealistic, highly aesthetic, refined, and professional.
🚫 No structural edits, only interior design transformations.`

func (c *OpenAIClient) getPrompt(style, roomtype string) string {
	return strings.ReplaceAll(strings.ReplaceAll(prompt, "[STYLE]", style), "[ROOM_TYPE]", roomtype)
}
