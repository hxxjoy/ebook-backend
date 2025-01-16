package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

var (
    nouns = []string{
        // 动物类
        "Panda", "Tiger", "Eagle", "Dolphin", "Fox", "Wolf", "Bear", "Lion", "Dragon", "Phoenix",
        "Rabbit", "Horse", "Deer", "Cat", "Dog", "Owl", "Hawk", "Falcon", "Leopard", "Whale",
        "Shark", "Penguin", "Koala", "Kangaroo", "Elephant", "Lynx", "Cheetah", "Panther", "Jaguar", "Raccoon",
        
        // 神话生物
        "Griffin", "Unicorn", "Pegasus", "Kraken", "Mermaid", "Sphinx", "Hydra", "Chimera", "Basilisk", "Wyrm",
        
        // 自然元素
        "Storm", "Thunder", "Lightning", "River", "Ocean", "Mountain", "Star", "Moon", "Sun", "Cloud",
        "Wind", "Rain", "Snow", "Frost", "Fire", "Flame", "Crystal", "Stone", "Forest", "Shadow",
        
        // 天体
        "Nova", "Comet", "Nebula", "Galaxy", "Meteor", "Aurora", "Orion", "Venus", "Mars", "Jupiter",
        
        // 植物
        "Oak", "Pine", "Maple", "Lotus", "Rose", "Lily", "Orchid", "Bamboo", "Cherry", "Willow",
        
        // 宝石
        "Ruby", "Sapphire", "Emerald", "Diamond", "Pearl", "Jade", "Amber", "Opal", "Topaz", "Garnet",
        
        // 抽象概念
        "Spirit", "Soul", "Dream", "Hope", "Destiny", "Legend", "Miracle", "Mystery", "Wonder", "Magic",
        
        // 其他生物
        "Raven", "Crow", "Serpent", "Cobra", "Viper", "Falcon", "Kirin", "Angel", "Knight", "Warrior",
        
        // 季节与时间
        "Spring", "Summer", "Autumn", "Winter", "Dawn", "Dusk", "Twilight", "Night", "Day", "Eclipse",
        
        // 地理特征
        "Peak", "Valley", "Canyon", "Cave", "Island", "Coast", "Cliff", "Reef", "Oasis", "Glacier",
    }
)

func GenerateNickname() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
    noun := nouns[r.Intn(len(nouns))]
    num := r.Intn(10000000)
    
    return fmt.Sprintf("%s%02d", noun, num)
}

func Md5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

