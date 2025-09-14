package gemini

import "strings"

const prompt = `Redesign this interior photograph into an aesthetically refined and super creative [STYLE] design.  
Keep the original room layout, geometry, perspective, and proportions exactly the same.  

🚫 Absolute restrictions (non-negotiable):  
- The architecture of the room must remain unchanged.  
- Windows, doors, and walls must stay in their exact original positions, size, and shape.  
- No new windows, doors, walls, or structural features may be created.  
- Do not move or alter ceiling height, room proportions, or camera perspective.  

Enhancements (allowed and encouraged):  
- Remove all clutter, construction items, old furniture, wires, or unnecessary objects.  
- Replace with elegant, imaginative, and well-arranged furniture and decorations that match the selected style.  
- Be **super creative** in furniture design, color palettes, lighting, textures, and decorative details — but only within the fixed structure of the original room.  
- Emphasize bold, unique, and visually striking design choices that are still harmonious.  
- Focus on balanced composition, visual storytelling, and aesthetic appeal.  
- Ensure the result looks like a professional interior design photoshoot for a high-end magazine.  
- Use natural lighting, realistic materials, and artistic staging for maximum impact.  

Guidelines:  
- Style: [STYLE] (Modern, Minimalist, Loft, Scandinavian, Japandi, Classic, etc.)  
- Room type: [ROOM_TYPE]  
- Output: photorealistic, highly aesthetic, refined, polished, and visually stunning.  
- Avoid dark, gloomy, or cluttered looks.  
- Absolutely no structural modifications, only interior design transformations.  
- The final image must look like a top-tier designer’s showcase, with the same architecture as the original room.  
`

func (c *Client) getPrompt(style, roomtype string) string {
	return strings.ReplaceAll(strings.ReplaceAll(prompt, "[STYLE]", style), "[ROOM_TYPE]", roomtype)
}
