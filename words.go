package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Difficulty represents word difficulty level
type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyMedium
	DifficultyHard
)

// DifficultyString returns a string representation of difficulty
func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "easy"
	case DifficultyMedium:
		return "medium"
	case DifficultyHard:
		return "hard"
	default:
		return "medium"
	}
}

// Easy words: short, common, simple (2-4 letters)
var easyWords = []string{
	"a", "an", "as", "at", "be", "by", "do", "go", "he", "if",
	"in", "is", "it", "me", "my", "no", "of", "on", "or", "so",
	"to", "up", "us", "we", "act", "add", "age", "ago", "aid", "air",
	"all", "and", "any", "are", "arm", "art", "ask", "bad", "bag", "bar",
	"bat", "bay", "bed", "bet", "bid", "big", "bin", "bit", "box", "boy",
	"bud", "bug", "bus", "but", "buy", "can", "car", "cat", "cop", "cow",
	"cry", "cup", "cut", "dad", "day", "die", "dig", "dim", "dog", "dot",
	"dry", "due", "dug", "ear", "eat", "egg", "ego", "elf", "end", "eye",
	"fan", "far", "fat", "fax", "fee", "few", "fig", "fit", "fix", "fly",
	"fog", "for", "fox", "fun", "fur", "gas", "get", "god", "got", "gum",
	"gun", "guy", "gym", "had", "has", "hat", "hay", "her", "hey", "hid",
	"him", "hip", "his", "hit", "hop", "hot", "how", "hub", "hug", "hut",
	"ice", "ill", "ink", "inn", "ion", "its", "jaw", "jay", "jet", "job",
	"jog", "joy", "key", "kid", "kin", "kit", "lab", "lad", "lap", "law",
	"lay", "leg", "let", "lid", "lie", "lip", "lit", "log", "lot", "low",
	"mad", "man", "map", "mat", "may", "men", "met", "mid", "mix", "mob",
	"mom", "mop", "mud", "mug", "nap", "net", "new", "nil", "nod", "nor",
	"not", "now", "nut", "oak", "odd", "off", "oil", "old", "one", "orb",
	"our", "out", "owe", "owl", "own", "pad", "pal", "pan", "par", "pay",
	"pea", "pen", "per", "pet", "pie", "pig", "pin", "pit", "pop", "pot",
	"pro", "pub", "pun", "put", "rag", "ram", "ran", "rap", "rat", "raw",
	"ray", "red", "rib", "rid", "rig", "rim", "rip", "rob", "rod", "row",
	"rub", "rug", "run", "sad", "sat", "saw", "say", "sea", "see", "set",
	"sew", "sex", "she", "shy", "sin", "sip", "sir", "sit", "six", "ski",
	"sky", "sly", "son", "spa", "spy", "sub", "sue", "sum", "sun", "tab",
	"tag", "tan", "tap", "tax", "tea", "ten", "the", "tie", "tin", "tip",
	"toe", "ton", "too", "top", "tow", "toy", "try", "tub", "tug", "two",
	"use", "van", "vet", "via", "war", "was", "wax", "way", "web", "wed",
	"wet", "who", "why", "win", "wit", "won", "woo", "wow", "yes", "yet",
	"you", "zoo",
}

// Medium words: common words (5-7 letters)
var mediumWords = []string{
	"about", "above", "abuse", "actor", "acute", "admit", "adopt", "adult", "after", "again",
	"agent", "agree", "ahead", "alarm", "album", "alert", "alike", "alive", "allow", "alone",
	"along", "alter", "among", "anger", "angle", "angry", "apart", "apple", "apply", "arena",
	"argue", "arise", "array", "aside", "asset", "audio", "audit", "avoid", "award", "aware",
	"badly", "baker", "bases", "basic", "basis", "beach", "began", "begin", "begun", "being",
	"below", "bench", "billy", "birth", "black", "blame", "blind", "block", "blood", "board",
	"brain", "brand", "bread", "break", "brick", "brief", "bring", "broad", "broke", "brown",
	"build", "built", "buyer", "cable", "calm", "canal", "carry", "catch", "cause", "chain",
	"chair", "chart", "chase", "cheap", "check", "chest", "chief", "child", "china", "chose",
	"civil", "claim", "class", "clean", "clear", "click", "clock", "close", "coach", "coast",
	"could", "count", "court", "cover", "craft", "crash", "cream", "crime", "cross", "crowd",
	"crown", "curve", "daily", "dance", "dated", "dealt", "death", "debut", "delay", "depth",
	"doing", "doubt", "dozen", "draft", "drama", "drawn", "dream", "dress", "drill", "drink",
	"drive", "drove", "dying", "early", "earth", "eight", "elite", "empty", "enemy", "enjoy",
	"enter", "entry", "equal", "error", "event", "every", "exact", "exist", "extra", "faith",
	"false", "fault", "fiber", "field", "fifth", "fifty", "fight", "final", "first", "fixed",
	"flash", "fleet", "floor", "fluid", "focus", "force", "forth", "forty", "forum", "found",
	"frame", "frank", "fraud", "fresh", "front", "fruit", "fully", "funny", "giant", "given",
	"glass", "globe", "going", "grace", "grade", "grand", "grant", "grass", "great", "green",
	"gross", "group", "grown", "guard", "guess", "guest", "guide", "happy", "harry", "heart",
	"heavy", "hello", "help", "hence", "horse", "hotel", "house", "human", "ideal", "image",
	"index", "inner", "input", "issue", "jeans", "joint", "judge", "juice", "knife", "knock",
	"known", "label", "large", "laser", "later", "laugh", "layer", "learn", "lease", "least",
	"leave", "legal", "level", "light", "limit", "links", "lives", "local", "logic", "loose",
	"lower", "lucky", "lunch", "lying", "magic", "major", "maker", "march", "marry", "match",
	"maybe", "mayor", "meant", "media", "metal", "might", "minor", "minus", "mixed", "model",
	"money", "month", "moral", "motor", "mount", "mouse", "mouth", "movie", "music", "needs",
	"never", "newly", "night", "noise", "north", "noted", "novel", "nurse", "occur", "ocean",
	"offer", "often", "order", "other", "ought", "paint", "panel", "paper", "party", "peace",
	"phase", "phone", "photo", "piece", "pilot", "pitch", "place", "plain", "plane", "plant",
	"plate", "plays", "plaza", "point", "pound", "power", "press", "price", "pride", "prime",
	"print", "prior", "prize", "proof", "proud", "prove", "queen", "quick", "quiet", "quite",
	"radio", "raise", "range", "rapid", "ratio", "reach", "ready", "refer", "right", "rival",
	"river", "robot", "roman", "rough", "round", "route", "royal", "rural", "scale", "scene",
	"scope", "score", "sense", "serve", "seven", "shall", "shape", "share", "sharp", "sheet",
	"shelf", "shell", "shift", "shirt", "shock", "shoot", "short", "shown", "sigh", "sight",
	"sign", "silent", "silly", "since", "skill", "sleep", "slice", "slide", "small", "smart",
	"smell", "smile", "smith", "smoke", "smooth", "snake", "solar", "solid", "solve", "sorry",
	"sound", "south", "space", "spare", "speak", "speed", "spend", "spent", "split", "spoke",
	"sport", "staff", "stage", "stake", "stand", "start", "state", "steam", "steel", "stick",
	"still", "stock", "stone", "stood", "store", "storm", "story", "strip", "stuck", "study",
	"stuff", "style", "sugar", "suite", "super", "sweet", "table", "taken", "taste", "taxes",
	"teach", "teeth", "terry", "texas", "thank", "theft", "their", "theme", "there", "these",
	"thick", "thing", "think", "third", "those", "three", "threw", "throw", "tight", "times",
	"tired", "title", "today", "topic", "total", "touch", "tough", "tower", "track", "trade",
	"train", "treat", "trend", "trial", "tried", "truck", "truly", "trust", "truth", "twice",
	"under", "undue", "union", "unity", "until", "upper", "upset", "urban", "usage", "usual",
	"value", "video", "virus", "visit", "vital", "voice", "waste", "watch", "water", "wheel",
	"where", "which", "while", "white", "whole", "whose", "woman", "women", "world", "worry",
	"worse", "worst", "worth", "would", "wound", "write", "wrong", "wrote", "yield", "young",
	"youth",
}

// Hard words: longer, complex, or unusual (8+ letters or difficult patterns)
var hardWords = []string{
	"absolute", "abstract", "academic", "accepted", "accident", "accuracy", "accurate", "achieved", "acquired", "activity",
	"actually", "addition", "adequate", "adjacent", "adjusted", "advanced", "affected", "afternoon", "agreement", "aircraft",
	"algorithm", "allocate", "allowing", "although", "aluminum", "analysis", "analyzing", "announced", "annually", "answered",
	"anticipate", "anything", "anywhere", "apparent", "apparently", "appearing", "appendix", "appetite", "applying", "approach",
	"approval", "approved", "architect", "argument", "arranged", "arrested", "arriving", "articles", "artificial", "aspirin",
	"assembly", "assuming", "assumption", "attempt", "attended", "attention", "attitude", "attorney", "attractive", "audience",
	"autumn", "averaging", "avoiding", "backbone", "background", "bacteria", "balanced", "barrier", "baseball", "bathroom",
	"becoming", "behavior", "believed", "belonging", "beneficial", "benefits", "betrayal", "birthday", "blessing", "blocking",
	"blossom", "boundary", "branches", "breaking", "breathing", "brilliant", "broccoli", "building", "bulletin", "business",
	"calendar", "campaign", "capacity", "carefully", "carrying", "category", "celebrate", "cemetery", "ceremony", "champion",
	"changing", "character", "charcoal", "chemical", "children", "choosing", "chronic", "cigarette", "cinnamon", "circular",
	"civilian", "clarify", "classical", "cleaning", "clearing", "clinical", "climbing", "clothing", "collapse", "colleague",
	"collected", "college", "colonial", "colorful", "combined", "commerce", "commercial", "committed", "committee", "commonly",
	"communist", "community", "companion", "comparing", "compelled", "competent", "competing", "complaint", "complete", "complex",
	"comply", "composed", "composer", "compound", "comprehend", "compromise", "computer", "conceal", "concede", "conceive",
	"concentrate", "concept", "concern", "concert", "conclude", "concrete", "condemn", "condition", "conduct", "conference",
	"confess", "confidence", "confident", "confine", "confirm", "conflict", "confuse", "congratulate", "congress", "connect",
	"conscious", "consensus", "consent", "conserve", "consider", "consist", "consistent", "constant", "constitute", "construct",
	"consult", "consume", "consumer", "contact", "contain", "contempt", "contend", "content", "contest", "context",
	"continue", "continuous", "contract", "contrary", "contrast", "contribute", "control", "controversy", "convenient", "conversation",
	"convert", "convince", "cooking", "cooperate", "coordinate", "coping", "copyright", "corner", "corporate", "correct",
	"corridor", "cottage", "council", "counselor", "counter", "countryside", "courageous", "courtesy", "coverage", "creating",
	"creative", "creature", "criminal", "crisis", "criterion", "critic", "critical", "criticism", "criticize", "crucial",
	"crushing", "cultural", "curiosity", "currency", "curriculum", "custody", "customary", "customer", "dangerous", "darkness",
	"database", "daughter", "daylight", "deadline", "dealing", "debate", "debris", "decade", "decide", "decision",
	"decisive", "declaration", "declare", "decline", "decorate", "decrease", "dedicate", "defeat", "defend", "defendant",
	"defense", "defensive", "deficit", "defined", "definitely", "definition", "delegate", "delicate", "delicious", "delight",
	"deliver", "delivery", "democracy", "democrat", "demonstrate", "department", "departure", "dependent", "depending", "depict",
	"depressed", "depression", "describe", "description", "desert", "deserve", "designer", "desirable", "desperate", "despite",
	"destination", "destroy", "destruction", "detective", "determine", "developing", "development", "devoted", "diabetes", "diagnose",
	"diagnosis", "diagram", "diamond", "dictate", "dictionary", "difference", "different", "difficult", "difficulty", "digital",
	"dignity", "dimension", "diminish", "dining", "diplomat", "directive", "directly", "director", "disability", "disabled",
	"disagree", "disappear", "disappoint", "disaster", "discipline", "disclose", "discount", "discourage", "discover", "discretion",
	"discriminate", "discuss", "disease", "disguise", "disgusting", "dishonest", "disk", "dismiss", "disorder", "display",
	"disposal", "dispute", "disrupt", "dissolve", "distance", "distant", "distinct", "distinction", "distinguish", "distort",
	"distract", "distress", "distribute", "district", "disturb", "diverse", "division", "divorce", "doctrine", "document",
	"domestic", "dominant", "dominate", "donation", "doorway", "double", "doubtful", "dramatic", "drawing", "driving",
	"dwelling", "dynamic", "earnest", "economic", "economist", "economy", "educate", "educator", "effective", "efficiency",
	"efficient", "effort", "elaborate", "election", "electoral", "electric", "elegant", "element", "elementary", "elevator",
	"eligible", "eliminate", "elite", "elsewhere", "embarrass", "embassy", "embrace", "emerge", "emergency", "emission",
	"emotion", "emotional", "emphasis", "emphasize", "empire", "employ", "employee", "employer", "employment", "empty",
	"enabling", "enclose", "encounter", "encourage", "endless", "endorse", "endure", "enemy", "energy", "enforce",
	"engage", "engine", "engineer", "enormous", "enough", "enrich", "enroll", "ensure", "enterprise", "entertain",
	"enthusiasm", "entirely", "entrance", "entries", "envelope", "environment", "environmental", "episode", "equality", "equation",
	"equipment", "equivalent", "error", "escape", "especially", "essential", "establish", "establishment", "estate", "estimate",
	"ethical", "ethics", "ethnic", "evaluate", "evaluation", "evening", "eventually", "everybody", "everyday", "everyone",
	"evidence", "evident", "evil", "evolution", "exactly", "examination", "examine", "example", "exceed", "excellent",
	"except", "exception", "excessive", "exchange", "excite", "excitement", "exciting", "exclude", "exclusion", "exclusive",
	"excuse", "execute", "execution", "executive", "exercise", "exhaust", "exhibit", "exhibition", "exist", "existence",
	"existing", "expand", "expansion", "expect", "expectation", "expense", "expensive", "experience", "experiment", "expert",
	"expertise", "explain", "explanation", "explode", "exploit", "exploration", "explore", "explosion", "expose", "exposure",
	"express", "expression", "extend", "extension", "extensive", "extent", "external", "extinct", "extra", "extraordinary",
	"extreme", "facility", "faculty", "failure", "fairly", "faith", "falling", "familiar", "family", "famous",
	"fantasy", "farmer", "fashion", "fasten", "fateful", "father", "fatigue", "fault", "favorite", "favorable",
	"feature", "federal", "fee", "feeding", "feeling", "fellow", "female", "feminist", "fence", "festival",
	"fever", "fiction", "field", "fierce", "fighting", "figure", "file", "filter", "final", "finance",
	"financial", "finding", "finger", "finish", "fire", "firmware", "first", "fiscal", "fishing", "fitness",
	"flame", "flash", "flavor", "flesh", "flexible", "flight", "floating", "flood", "floor", "flourish",
	"flower", "fluid", "focus", "fog", "fold", "folk", "following", "foot", "forbid", "force",
	"foreign", "forest", "forever", "forget", "forgive", "formal", "formation", "former", "formula", "forth",
	"fortunate", "fortune", "forum", "forward", "foster", "found", "foundation", "founder", "fraction", "fragile",
	"fragment", "frame", "framework", "franchise", "frank", "fraud", "freedom", "freeway", "freeze", "freight",
	"frequency", "frequent", "fresh", "friction", "friendly", "friendship", "frighten", "frog", "frontier", "frown",
	"fruit", "frustrate", "fuel", "fulfill", "full", "fully", "fun", "function", "fund", "fundamental",
	"funding", "funeral", "funny", "furniture", "furthermore", "future", "galaxy", "gallery", "gambling", "gaming",
	"garage", "garden", "garlic", "garment", "gasoline", "gather", "gathering", "gauge", "gaze", "gear",
	"gender", "general", "generate", "generation", "generous", "genius", "genre", "gentle", "genuine", "gesture",
	"getting", "ghost", "giant", "gifted", "girl", "girlfriend", "give", "given", "glad", "glance",
	"glass", "global", "glory", "glove", "go", "goal", "goat", "gold", "golden", "golf",
	"goodbye", "goodness", "goods", "govern", "government", "governor", "grace", "grade", "gradual", "graduate",
	"graduation", "grain", "grand", "grandfather", "grandmother", "grant", "grape", "grass", "grateful", "gravity",
	"great", "greatest", "green", "greenhouse", "greet", "grief", "grin", "grind", "grip", "grocery",
	"ground", "group", "grow", "growing", "growth", "guarantee", "guard", "guess", "guest", "guide",
	"guideline", "guilt", "guilty", "guitar", "gun", "gut", "guy", "gym", "habit", "habitat",
	"hail", "hair", "half", "hall", "hallway", "hand", "handful", "handle", "handsome", "hang",
	"happen", "happily", "happiness", "happy", "harbor", "hard", "hardly", "hardware", "harm", "harmful",
	"harmony", "harsh", "harvest", "hat", "hate", "haul", "have", "hazard", "head", "headache",
	"headline", "headquarters", "heal", "health", "healthy", "hear", "hearing", "heart", "heat", "heaven",
	"heavily", "heavy", "heel", "height", "heir", "helicopter", "hell", "hello", "help", "helpful",
	"helpless", "hence", "herb", "heritage", "hero", "heroic", "herself", "hesitate", "hidden", "hide",
	"hierarchy", "high", "highly", "highway", "hike", "hill", "himself", "hint", "hip", "hire",
	"historian", "historic", "historical", "history", "hit", "hobby", "hold", "holding", "hole", "holiday",
	"holy", "home", "homeless", "homework", "honest", "honesty", "honey", "honor", "hook", "hope",
	"hopeful", "hopefully", "hopeless", "horizon", "hormone", "horn", "horrible", "horror", "horse", "hospital",
	"host", "hostile", "hot", "hotel", "hour", "house", "household", "housing", "however", "huge",
	"human", "humanitarian", "humanity", "humble", "humor", "hundred", "hunger", "hungry", "hunt", "hunter",
	"hunting", "hurricane", "hurry", "hurt", "husband", "hypothesis", "ice", "icon", "idea", "ideal",
	"identical", "identification", "identify", "identity", "ideology", "if", "ignorance", "ignore", "ill", "illegal",
	"illness", "illusion", "illustrate", "image", "imagination", "imagine", "immediate", "immediately", "immigrant", "immigration",
	"immune", "impact", "implement", "implementation", "implication", "imply", "import", "importance", "important", "impose",
	"impossible", "impress", "impression", "impressive", "improve", "improvement", "incentive", "incident", "include", "including",
	"income", "incorporate", "increase", "increasing", "increasingly", "incredible", "indeed", "independence", "independent", "index",
	"indicate", "indication", "indicator", "individual", "industrial", "industry", "inevitable", "infant", "infection", "inflation",
	"influence", "influential", "inform", "informal", "information", "infrastructure", "ingredient", "inhabitant", "inherent", "inherit",
	"inhibit", "initial", "initially", "initiate", "initiative", "injure", "injury", "inmate", "inner", "innocent",
	"innovation", "innovative", "input", "inquiry", "insect", "insert", "inside", "insight", "insist", "inspect",
	"inspection", "inspector", "inspire", "install", "installation", "instance", "instant", "instead", "instinct", "institute",
	"institution", "institutional", "instruct", "instruction", "instructor", "instrument", "insufficient", "insurance", "intellectual", "intelligence",
	"intelligent", "intend", "intense", "intensity", "intent", "intention", "interact", "interaction", "interest", "interested",
	"interesting", "interfere", "interior", "internal", "international", "internet", "interpret", "interpretation", "interrupt", "interval",
	"intervention", "interview", "intimate", "into", "introduce", "introduction", "invade", "invasion", "invent", "invention",
	"inventory", "invest", "investigate", "investigation", "investigator", "investment", "investor", "invisible", "invitation", "invite",
	"involve", "involved", "involvement", "iron", "irony", "island", "isolate", "isolated", "isolation", "issue",
	"item", "its", "itself", "jacket", "jail", "jar", "jaw", "jazz", "jealous", "jeans",
	"jet", "jewelry", "job", "join", "joint", "joke", "journal", "journalism", "journalist", "journey",
	"joy", "judge", "judgment", "judicial", "juice", "jump", "jungle", "junior", "jury", "just",
	"justice", "justify", "keen", "keep", "keeper", "kettle", "key", "keyboard", "kick", "kid",
	"kill", "killer", "killing", "kind", "king", "kingdom", "kiss", "kit", "kitchen", "knee",
	"kneel", "knife", "knit", "knock", "knot", "know", "knowledge", "label", "labor", "laboratory",
	"laborer", "lack", "ladder", "lady", "lake", "lamb", "lamp", "land", "landing", "landlord",
	"landscape", "lane", "language", "large", "largely", "last", "late", "lately", "later", "latest",
	"latter", "laugh", "launch", "law", "lawn", "lawsuit", "lawyer", "lay", "layer", "layout",
	"lazy", "lead", "leader", "leadership", "leading", "leaf", "league", "lean", "leap", "learn",
	"learning", "lease", "least", "leather", "leave", "lecture", "left", "leg", "legacy", "legal",
	"legend", "legislation", "legislative", "legislator", "legislature", "legitimate", "lemon", "lend", "length", "lens",
	"less", "lesser", "lesson", "let", "letter", "level", "liability", "liberal", "liberty", "library",
	"license", "lie", "life", "lifestyle", "lifetime", "lift", "light", "lightning", "like", "likely",
	"likewise", "limb", "limit", "limitation", "limited", "line", "link", "lion", "lip", "liquid",
	"list", "listen", "listener", "literally", "literary", "literature", "little", "live", "liver", "living",
	"load", "loan", "lobby", "local", "locale", "locate", "location", "lock", "log", "logic",
	"logical", "lonely", "long", "longtime", "look", "loose", "lose", "loss", "lost", "lot",
	"lots", "loud", "love", "lovely", "lover", "low", "lower", "loyal", "loyalty", "luck",
	"lucky", "lunch", "lung", "luxury", "lying", "machine", "machinery", "mad", "magazine", "magic",
	"magnetic", "magnificent", "mail", "main", "mainly", "mainstream", "maintain", "maintenance", "major", "majority",
	"make", "maker", "makeup", "male", "mall", "man", "manage", "management", "manager", "managing",
	"mandate", "manipulate", "manner", "manufacturer", "manufacturing", "many", "map", "margin", "marine", "mark",
	"market", "marketing", "marketplace", "marriage", "married", "marry", "marsh", "marvelous", "mask", "mass",
	"massive", "master", "match", "mate", "material", "math", "mathematics", "matter", "mature", "maximum",
	"may", "maybe", "mayor", "me", "meal", "mean", "meaning", "meaningful", "meantime", "meanwhile",
	"measure", "measurement", "meat", "mechanic", "mechanical", "mechanism", "medal", "media", "medical", "medication",
	"medicine", "medium", "meet", "meeting", "melt", "member", "membership", "memory", "mental", "mentally",
	"mention", "mentor", "menu", "merchant", "mere", "merely", "merit", "mess", "message", "metal",
	"metaphor", "meter", "method", "methodology", "middle", "midnight", "midst", "might", "migration", "mild",
	"military", "milk", "mill", "million", "mind", "mine", "mineral", "minimal", "minimize", "minimum",
	"minister", "ministry", "minor", "minority", "minute", "miracle", "mirror", "miss", "missile", "mission",
	"missionary", "mistake", "mix", "mixture", "moan", "mobile", "mode", "model", "moderate", "modern",
	"modest", "modify", "molecule", "mom", "moment", "momentum", "money", "monitor", "monkey", "monster",
	"month", "monument", "mood", "moon", "moral", "morality", "more", "moreover", "morning", "mortgage",
	"most", "mostly", "mother", "motion", "motivate", "motivation", "motive", "motor", "mount", "mountain",
	"mouse", "mouth", "move", "movement", "movie", "much", "mud", "multiple", "murder", "muscle",
	"museum", "music", "musical", "musician", "must", "mutter", "mutual", "my", "myself", "mysterious",
	"mystery", "myth", "nail", "naked", "name", "narrative", "narrow", "nation", "national", "nationwide",
	"native", "natural", "naturally", "nature", "naval", "navigation", "near", "nearby", "nearly", "neat",
	"necessarily", "necessary", "necessity", "neck", "need", "needed", "needing", "negative", "neglect", "negotiate",
	"negotiation", "neighbor", "neighborhood", "neither", "nerve", "nervous", "nest", "net", "network", "neutral",
	"never", "nevertheless", "new", "newly", "news", "newspaper", "next", "nice", "night", "nightmare",
	"nine", "no", "nobody", "nod", "noise", "nomination", "nominee", "none", "nonetheless", "nonprofit",
	"noon", "nor", "normal", "normally", "north", "northern", "nose", "not", "note", "notebook",
	"nothing", "notice", "notion", "novel", "now", "nowhere", "nuclear", "number", "numerous", "nurse",
	"nursery", "nut", "nutrition", "oak", "object", "objection", "objective", "obligation", "observation", "observe",
	"observer", "obstacle", "obtain", "obvious", "obviously", "occasion", "occasional", "occasionally", "occupation", "occupy",
	"occur", "ocean", "odd", "odds", "off", "offense", "offensive", "offer", "offering", "office",
	"officer", "official", "officially", "often", "oil", "okay", "old", "ongoing", "online", "only",
	"onto", "open", "opening", "openly", "opera", "operate", "operating", "operation", "operator", "opinion",
	"opponent", "opportunity", "oppose", "opposite", "opposition", "opt", "optical", "optimistic", "option", "optional",
	"oral", "orange", "orbit", "orchestra", "order", "ordinary", "organ", "organic", "organization", "organizational",
	"organize", "organized", "organizer", "orientation", "origin", "original", "originally", "originate", "other", "otherwise",
	"ought", "ourselves", "out", "outcome", "outdoor", "outer", "outfit", "outlet", "outline", "output",
	"outrage", "outside", "outsider", "outstanding", "oven", "over", "overall", "overcome", "overlook", "oversee",
	"overwhelm", "overwhelming", "owe", "own", "owner", "ownership", "oxygen", "pace", "pack", "package",
	"pad", "page", "pain", "painful", "paint", "painter", "painting", "pair", "pal", "palace",
	"pale", "palm", "pan", "panel", "pants", "paper", "parade", "parent", "parental", "parish",
	"park", "parking", "part", "partial", "partially", "participant", "participate", "participation", "particle", "particular",
	"particularly", "partly", "partner", "partnership", "party", "pass", "passage", "passenger", "passing", "passion",
	"past", "patch", "patent", "path", "patience", "patient", "patrol", "patron", "pattern", "pause",
	"pay", "payment", "peace", "peaceful", "peak", "peer", "pen", "penalty", "pencil", "people",
	"pepper", "per", "perceive", "percentage", "perception", "perfect", "perfectly", "perform", "performance", "performer",
	"perhaps", "period", "permanent", "permission", "permit", "person", "personal", "personality", "personally", "personnel",
	"perspective", "persuade", "pet", "phase", "phenomenon", "philosophy", "phone", "photo", "photograph", "photographer",
	"photography", "phrase", "physical", "physically", "physician", "physics", "piano", "pick", "picture", "pie",
	"piece", "pig", "pile", "pill", "pilot", "pine", "pink", "pioneer", "pipe", "pitch",
	"pity", "place", "plain", "plan", "plane", "planet", "planning", "plant", "plastic", "plate",
	"platform", "play", "player", "plea", "pleasant", "please", "pleased", "pleasure", "plenty", "plot",
	"plug", "plunge", "plus", "pocket", "poem", "poet", "poetry", "point", "pole", "police",
	"policeman", "policy", "political", "politically", "politician", "politics", "poll", "pollution", "pond", "pool",
	"poor", "pop", "popular", "popularity", "population", "porch", "pork", "port", "portable", "portfolio",
	"portion", "portrait", "portray", "pose", "position", "positive", "possess", "possession", "possibility", "possible",
	"possibly", "post", "poster", "pot", "potato", "potential", "potentially", "pound", "pour", "poverty",
	"powder", "power", "powerful", "practical", "practically", "practice", "pray", "prayer", "preach", "precious",
	"precise", "precisely", "predator", "predict", "predictable", "prediction", "prefer", "preference", "pregnancy", "pregnant",
	"preliminary", "premise", "premium", "preparation", "prepare", "prescription", "presence", "present", "presentation", "preserve",
	"presidency", "president", "presidential", "press", "pressure", "pretend", "pretty", "prevail", "prevent", "prevention",
	"previous", "previously", "price", "pride", "priest", "primarily", "primary", "prime", "principal", "principle",
	"print", "prior", "priority", "prison", "prisoner", "privacy", "private", "privately", "privilege", "prize",
	"pro", "probably", "problem", "procedure", "proceed", "process", "processing", "processor", "proclaim", "produce",
	"producer", "product", "production", "productive", "productivity", "profession", "professional", "professor", "profile", "profit",
	"profound", "program", "programming", "progress", "progressive", "prohibit", "project", "prominent", "promise", "promising",
	"promote", "promotion", "prompt", "proof", "proper", "properly", "property", "proportion", "proposal", "propose",
	"proposed", "prosecution", "prospect", "protect", "protection", "protective", "protein", "protest", "protocol", "proud",
	"prove", "provide", "provided", "provider", "province", "provincial", "provision", "psychological", "psychologist", "psychology",
	"public", "publication", "publicity", "publicly", "publish", "publisher", "pull", "pulse", "pump", "punch",
	"punish", "punishment", "purchase", "pure", "purpose", "pursue", "pursuit", "push", "put", "puzzle",
	"qualify", "quality", "quantity", "quarter", "queen", "quest", "question", "quick", "quickly", "quiet",
	"quietly", "quit", "quite", "quote", "race", "racial", "racing", "racism", "racist", "radar",
	"radiation", "radical", "radio", "rage", "rail", "railroad", "rain", "raise", "rally", "ranch",
	"random", "range", "rank", "rapid", "rapidly", "rare", "rarely", "rat", "rate", "rather",
	"rating", "ratio", "rational", "raw", "ray", "reach", "react", "reaction", "read", "reader",
	"reading", "ready", "real", "realistic", "reality", "realize", "really", "reason", "reasonable", "reasonably",
	"reasoning", "rebel", "rebellion", "rebuild", "recall", "receipt", "receive", "receiver", "recent", "recently",
	"reception", "recession", "recipe", "recipient", "reckon", "recognition", "recognize", "recommend", "recommendation", "record",
	"recording", "recover", "recovery", "recruit", "recruitment", "recycling", "red", "reduce", "reduction", "refer",
	"reference", "reflect", "reflection", "reform", "refugee", "refuse", "regard", "regarding", "regardless", "regime",
	"region", "regional", "register", "registration", "regret", "regular", "regularly", "regulate", "regulation", "regulatory",
	"rehabilitation", "reinforce", "reject", "relate", "relation", "relationship", "relative", "relatively", "relax", "release",
	"relevant", "reliable", "relief", "relieve", "religion", "religious", "reluctant", "rely", "remain", "remaining",
	"remark", "remarkable", "remedy", "remember", "remind", "reminder", "remote", "removal", "remove", "render",
	"renew", "renowned", "rent", "rental", "repair", "repeat", "repeatedly", "replace", "replacement", "report",
	"reporter", "represent", "representation", "representative", "reproduce", "republic", "reputation", "request", "require", "requirement",
	"rescue", "research", "researcher", "resemble", "reservation", "reserve", "residence", "resident", "residential", "resign",
	"resist", "resistance", "resolution", "resolve", "resort", "resource", "respect", "respectively", "respond", "respondent",
	"response", "responsibility", "responsible", "rest", "restaurant", "restore", "restrict", "restriction", "result", "resume",
	"retail", "retain", "retire", "retired", "retirement", "retreat", "retrieve", "return", "reveal", "revelation",
	"revenge", "revenue", "reverse", "review", "revise", "revision", "revival", "revive", "revolution", "revolutionary",
	"reward", "rhetoric", "rhythm", "rice", "rich", "rid", "ride", "rider", "ridge", "ridiculous",
	"rifle", "right", "rim", "ring", "riot", "rip", "rise", "rising", "risk", "risky",
	"ritual", "rival", "river", "road", "robot", "robust", "rock", "rocket", "rod", "role",
	"roll", "rolling", "romance", "romantic", "roof", "room", "root", "rope", "rose", "rough",
	"roughly", "round", "route", "routine", "row", "royal", "rub", "ruin", "rule", "ruling",
	"rumor", "run", "running", "rural", "rush", "sack", "sacred", "sacrifice", "sad", "safe",
	"safety", "sail", "sailing", "saint", "sake", "salad", "salary", "sale", "sales", "salmon",
	"salon", "salt", "same", "sample", "sanction", "sand", "sandwich", "satellite", "satisfaction", "satisfy",
	"sauce", "save", "saving", "say", "scale", "scan", "scandal", "scared", "scary", "scatter",
	"scenario", "scene", "scent", "schedule", "scheme", "scholar", "scholarship", "school", "science", "scientific",
	"scientist", "scope", "score", "scratch", "scream", "screen", "screening", "script", "sculpture", "sea",
	"search", "season", "seat", "second", "secondary", "secret", "secretary", "section", "sector", "secular",
	"secure", "security", "see", "seed", "seek", "seem", "seemingly", "segment", "seize", "seldom",
	"select", "selection", "selective", "self", "sell", "seller", "seminar", "senator", "send", "senior",
	"sensation", "sense", "sensible", "sensitive", "sentence", "sentiment", "separate", "sequence", "series", "serious",
	"seriously", "servant", "serve", "service", "serving", "session", "set", "setting", "settle", "settlement",
	"seven", "seventeen", "seventh", "seventy", "several", "severe", "severely", "severity", "sex", "sexual",
	"sexuality", "sexy", "shade", "shadow", "shake", "shall", "shallow", "shame", "shape", "share",
	"shared", "shark", "sharp", "sharply", "shatter", "she", "shed", "sheep", "sheer", "sheet",
	"shelf", "shell", "shelter", "shift", "shining", "ship", "shipping", "shirt", "shock", "shoe",
	"shoot", "shooting", "shop", "shopping", "shore", "short", "shortage", "shortly", "shorts", "shot",
	"should", "shoulder", "shout", "show", "shower", "shrug", "shut", "shy", "sick", "sickness",
	"side", "sidewalk", "sigh", "sight", "sign", "signal", "signature", "significance", "significant", "significantly",
	"silence", "silent", "silk", "silly", "silver", "similar", "similarity", "similarly", "simple", "simply",
	"sin", "since", "sincere", "sing", "singer", "single", "sink", "sir", "sister", "sit",
	"site", "situation", "six", "sixteen", "sixth", "sixty", "size", "ski", "skill", "skilled",
	"skin", "skip", "skirt", "skull", "sky", "slap", "slave", "slavery", "sleep", "slice",
	"slide", "slight", "slightly", "slip", "slogan", "slope", "slot", "slow", "slowly", "small",
	"smart", "smash", "smell", "smile", "smoke", "smooth", "snap", "snow", "so", "soak",
	"soap", "soar", "soccer", "social", "socially", "society", "soft", "software", "soil", "solar",
	"soldier", "sole", "solely", "solicit", "solid", "solo", "solution", "solve", "some", "somebody",
	"somehow", "someone", "something", "sometime", "sometimes", "somewhat", "somewhere", "son", "song", "soon",
	"sophisticated", "sorry", "sort", "soul", "sound", "soup", "source", "south", "southern", "space",
	"spare", "spark", "speak", "speaker", "special", "specialist", "specialize", "specialty", "species", "specific",
	"specifically", "specify", "specimen", "spectacle", "spectacular", "spectrum", "speculate", "speculation", "speech", "speed",
	"spell", "spelling", "spend", "spending", "sphere", "spill", "spin", "spine", "spirit", "spiritual",
	"spit", "spite", "split", "spoil", "spoken", "spokesman", "sponge", "sponsor", "spontaneous", "spoon",
	"sport", "sports", "spot", "spouse", "spray", "spread", "spring", "squad", "square", "squeeze",
	"stable", "stadium", "staff", "stage", "stair", "stake", "stall", "stamp", "stance", "stand",
	"standard", "standing", "star", "stare", "start", "starter", "starting", "starve", "state", "statement",
	"station", "statistical", "statistics", "statue", "status", "statute", "stay", "steady", "steal", "steam",
	"steel", "steep", "steer", "stem", "step", "stereotype", "stick", "stiff", "still", "stimulate",
	"stimulus", "stir", "stock", "stomach", "stone", "stop", "storage", "store", "storm", "story",
	"straight", "strain", "strange", "stranger", "strategic", "strategy", "stream", "street", "strength", "strengthen",
	"stress", "stretch", "strict", "strictly", "strike", "striking", "string", "strip", "strive", "stroke",
	"strong", "strongly", "structural", "structure", "struggle", "student", "studio", "study", "stuff", "stumble",
	"stupid", "style", "subject", "subjective", "submit", "subsequent", "subsequently", "subsidy", "substance", "substantial",
	"substantially", "substitute", "subtle", "suburb", "suburban", "succeed", "success", "successful", "successfully", "succession",
	"successive", "successor", "such", "sudden", "suddenly", "sue", "suffer", "suffering", "sufficient", "sufficiently",
	"sugar", "suggest", "suggestion", "suicide", "suit", "suitable", "suite", "sum", "summarize", "summary",
	"summer", "summit", "sun", "super", "superb", "superior", "supervisor", "supplement", "supply", "support",
	"supporter", "supportive", "suppose", "supposed", "supposedly", "suppress", "supreme", "sure", "surely", "surface",
	"surge", "surgeon", "surgery", "surgical", "surplus", "surprise", "surprised", "surprising", "surprisingly", "surrender",
	"surround", "surrounding", "survey", "survival", "survive", "survivor", "suspect", "suspend", "suspicion", "suspicious",
	"sustain", "sustainable", "swallow", "swear", "sweat", "sweep", "sweet", "swell", "swift", "swim",
	"swimming", "swing", "switch", "sword", "symbol", "symbolic", "sympathy", "symptom", "syndrome", "system",
	"systematic", "table", "tactic", "tail", "take", "tale", "talent", "talented", "talk", "tall",
	"tank", "tap", "tape", "target", "task", "taste", "tax", "taxation", "taxpayer", "tea",
	"teach", "teacher", "teaching", "team", "teammate", "tear", "teaspoon", "technical", "technician", "technique",
	"technological", "technology", "teen", "teenage", "teenager", "telephone", "telescope", "television", "tell", "temperature",
	"temple", "temporarily", "temporary", "ten", "tend", "tendency", "tender", "tennis", "tension", "tent",
	"term", "terminal", "terminate", "terms", "terrain", "terrible", "terribly", "terrific", "territory", "terror",
	"terrorism", "terrorist", "test", "testify", "testimony", "testing", "text", "textbook", "texture", "than",
	"thank", "thanks", "that", "the", "theater", "their", "them", "theme", "themselves", "then",
	"theological", "theology", "theoretical", "theory", "therapist", "therapy", "there", "thereby", "therefore", "these",
	"they", "thick", "thigh", "thin", "thing", "think", "thinking", "third", "thirty", "this",
	"those", "though", "thought", "thousand", "thread", "threat", "threaten", "three", "threshold", "thrilled",
	"thrive", "throat", "through", "throughout", "throw", "thumb", "thus", "ticket", "tide", "tie",
	"tight", "tighten", "tightly", "tile", "till", "timber", "time", "timing", "tin", "tiny",
	"tip", "tire", "tired", "tissue", "title", "to", "tobacco", "today", "toe", "together",
	"toilet", "tolerance", "tolerate", "toll", "tomato", "tomorrow", "ton", "tone", "tongue", "tonight",
	"too", "tool", "tooth", "top", "topic", "toss", "total", "totally", "touch", "tough",
	"tour", "tourism", "tourist", "tournament", "toward", "towards", "tower", "town", "toxic", "trace",
	"track", "trade", "trading", "tradition", "traditional", "traditionally", "traffic", "tragedy", "tragic", "trail",
	"trailer", "train", "training", "trait", "transaction", "transcript", "transfer", "transform", "transformation", "transit",
	"transition", "translate", "translation", "transmission", "transmit", "transport", "transportation", "trap", "trash", "trauma",
	"travel", "traveler", "tray", "treasure", "treat", "treatment", "treaty", "tree", "tremendous", "trend",
	"trial", "tribal", "tribe", "trick", "trigger", "trim", "trip", "triumph", "troop", "tropical",
	"trouble", "troubled", "truck", "true", "truly", "trunk", "trust", "trustee", "truth", "try",
	"tube", "tuck", "tuition", "tumble", "tune", "tunnel", "turn", "turnout", "turnover", "tutor",
	"twelve", "twentieth", "twenty", "twice", "twin", "twist", "type", "typical", "typically", "tyranny",
	"ugly", "ultimate", "ultimately", "unable", "unacceptable", "uncertain", "uncertainty", "uncle", "uncomfortable", "unconscious",
	"under", "undergo", "undergraduate", "underlying", "undermine", "understand", "understanding", "undertake", "unemployed", "unemployment",
	"unexpected", "unfair", "unfold", "unfortunate", "unfortunately", "unhappy", "uniform", "unify", "union", "unique",
	"unit", "unite", "united", "unity", "universal", "universe", "university", "unknown", "unless", "unlike",
	"unlikely", "unnecessary", "unprecedented", "until", "unusual", "unusually", "up", "update", "upgrade", "upon",
	"upper", "upset", "upstairs", "urban", "urge", "urgent", "us", "usage", "use", "used",
	"useful", "useless", "user", "usual", "usually", "utility", "utilize", "utterly", "vacation", "vaccine",
	"vacuum", "vague", "valid", "validity", "valley", "valuable", "value", "van", "vanish", "variable",
	"variation", "variety", "various", "vary", "vast", "vegetable", "vehicle", "venture", "verbal", "verdict",
	"version", "versus", "vertical", "very", "vessel", "veteran", "viable", "vibrant", "vice", "victim",
	"victory", "video", "view", "viewer", "viewpoint", "village", "violate", "violation", "violence", "violent",
	"virtual", "virtually", "virtue", "virus", "visible", "vision", "visit", "visitor", "visual", "vital",
	"vitamin", "vivid", "vocabulary", "vocal", "voice", "volume", "voluntary", "volunteer", "vote", "voter",
	"voting", "vow", "vulnerable", "wage", "wagon", "waist", "wait", "wake", "walk", "wall",
	"wander", "want", "war", "warm", "warmth", "warn", "warning", "warrant", "warrior", "wash",
	"waste", "watch", "water", "wave", "way", "we", "weak", "weaken", "weakness", "wealth",
	"wealthy", "weapon", "wear", "weather", "weave", "wedding", "week", "weekend", "weekly", "weigh",
	"weight", "weird", "welcome", "welfare", "well", "wellbeing", "west", "western", "wet", "whale",
	"what", "whatever", "whatsoever", "wheat", "wheel", "when", "whenever", "where", "whereas", "wherever",
	"whether", "which", "while", "whip", "whisper", "whistle", "white", "who", "whoever", "whole",
	"wholly", "whom", "whose", "why", "wicked", "wide", "widely", "widen", "widespread", "widow",
	"width", "wife", "wild", "wilderness", "wildlife", "will", "willing", "willingness", "win", "wind",
	"window", "wine", "wing", "winner", "winter", "wipe", "wire", "wisdom", "wise", "wish",
	"wit", "withdraw", "withdrawal", "within", "without", "witness", "wolf", "woman", "wonder", "wonderful",
	"wood", "wooden", "wool", "word", "work", "worker", "working", "workplace", "works", "workshop",
	"world", "worldwide", "worm", "worried", "worry", "worse", "worship", "worst", "worth", "worthwhile",
	"worthy", "would", "wound", "wow", "wrap", "wrath", "wrist", "write", "writer", "writing",
	"written", "wrong", "wrote", "yard", "yeah", "year", "yell", "yellow", "yes", "yesterday",
	"yet", "yield", "young", "youngster", "your", "yours", "yourself", "youth", "zero", "zone",
}

// DifficultyList returns the word list for a given difficulty
func DifficultyList(d Difficulty) []string {
	switch d {
	case DifficultyEasy:
		return easyWords
	case DifficultyHard:
		return hardWords
	default:
		return mediumWords
	}
}

// getRandomWords returns n random words from the word list for given difficulty
// Ensures no word appears consecutively for better typing flow
func getRandomWords(n int, difficulty Difficulty) []string {
	if n <= 0 {
		return []string{}
	}

	wordList := DifficultyList(difficulty)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		word := wordList[rand.Intn(len(wordList))]
		// Prevent consecutive duplicates
		if i > 0 && word == words[i-1] && len(wordList) > 1 {
			// Pick a different word
			for word == words[i-1] {
				word = wordList[rand.Intn(len(wordList))]
			}
		}
		words[i] = word
	}
	return words
}

// getDifficultyFromString parses difficulty from string
func getDifficultyFromString(s string) Difficulty {
	switch s {
	case "easy":
		return DifficultyEasy
	case "hard":
		return DifficultyHard
	default:
		return DifficultyMedium
	}
}

// Punctuation marks for advanced typing practice
var punctuationMarks = []string{
	",", ".", ";", ":", "!", "?", "-", "_", "'", "\"",
	"(", ")", "[", "]", "{", "}", "/", "\\", "@", "#",
	"$", "%", "^", "&", "*", "+", "=", "<", ">", "|",
}

// Common punctuation combinations
var punctuationCombos = []string{
	", ", ". ", "; ", ": ", "! ", "? ", " - ", "'s", "n't",
	"'re", "'ll", "'d", "'ve", "'m", "..", "...", "!!", "??",
}

// Numbers for numeric typing practice
var numberList = []string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19",
	"20", "25", "30", "40", "50", "60", "75", "80", "90", "100",
	"123", "456", "789", "1000", "2024", "3.14", "1.5", "2.0", "0.5",
}

// Number with punctuation combinations
var numberPunctuationCombos = []string{
	"1,000", "10,000", "100,000", "1,000,000", "3.14", "$100", "50%",
	"(2024)", "v1.0", "2.5", "-5", "+10", "1st", "2nd", "3rd", "4th",
	"5:30", "12:00", "24/7", "9-5", "20/20", "1:1", "2x", "3D",
}

// WordComplexity represents additional complexity options
type WordComplexity int

const (
	ComplexityNormal WordComplexity = iota
	ComplexityPunctuation
	ComplexityNumbers
	ComplexityFull
)

// ComplexityString returns a string representation
func (c WordComplexity) String() string {
	switch c {
	case ComplexityPunctuation:
		return "punctuation"
	case ComplexityNumbers:
		return "numbers"
	case ComplexityFull:
		return "full"
	default:
		return "normal"
	}
}

// getRandomWordsWithComplexity returns words with optional punctuation/numbers
func getRandomWordsWithComplexity(n int, difficulty Difficulty, complexity WordComplexity) []string {
	if n <= 0 {
		return []string{}
	}

	baseWords := getRandomWords(n, difficulty)
	if complexity == ComplexityNormal {
		return baseWords
	}

	words := make([]string, n)

	for i := 0; i < n; i++ {
		word := baseWords[i]

		switch complexity {
		case ComplexityPunctuation:
			// 30% chance to add punctuation to a word
			if rand.Float32() < 0.3 {
				word = addPunctuationToWord(word)
			}
		case ComplexityNumbers:
			// 20% chance to replace word with number, 10% to add number to word
			r := rand.Float32()
			if r < 0.2 {
				word = numberList[rand.Intn(len(numberList))]
			} else if r < 0.3 {
				word = word + numberList[rand.Intn(len(numberList))]
			}
		case ComplexityFull:
			// Mix of punctuation and numbers
			r := rand.Float32()
			if r < 0.25 {
				// Replace with number
				word = numberList[rand.Intn(len(numberList))]
			} else if r < 0.45 {
				// Add punctuation
				word = addPunctuationToWord(word)
			} else if r < 0.55 {
				// Number-punctuation combo
				word = numberPunctuationCombos[rand.Intn(len(numberPunctuationCombos))]
			}
		}

		words[i] = word
	}

	return words
}

// addPunctuationToWord adds punctuation to a word
func addPunctuationToWord(word string) string {
	// Different punctuation strategies
	switch rand.Intn(5) {
	case 0:
		// Add trailing punctuation
		return word + punctuationMarks[rand.Intn(len(punctuationMarks))]
	case 1:
		// Add leading punctuation
		return punctuationMarks[rand.Intn(len(punctuationMarks))] + word
	case 2:
		// Wrap in punctuation
		p1 := punctuationMarks[rand.Intn(len(punctuationMarks))]
		p2 := punctuationMarks[rand.Intn(len(punctuationMarks))]
		return p1 + word + p2
	case 3:
		// Add apostrophe combo
		return word + punctuationCombos[rand.Intn(6)+8] // 's, n't, 're, 'll, 'd, 've, 'm
	default:
		// Add internal punctuation
		if len(word) > 2 {
			mid := len(word) / 2
			return word[:mid] + punctuationMarks[rand.Intn(len(punctuationMarks))] + word[mid:]
		}
		return word + punctuationMarks[rand.Intn(len(punctuationMarks))]
	}
}
