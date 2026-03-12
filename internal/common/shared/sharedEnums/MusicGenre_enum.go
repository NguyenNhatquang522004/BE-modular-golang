package sharedEnums

//go:generate enumer -type=MusicGenre -json -transform=snake -trimprefix=Genre
type MusicGenre int

const (
	GenrePop        MusicGenre = iota // 'pop'
	GenreRock                         // 'rock'
	GenreHipHop                       // 'hip_hop'
	GenreRNB                          // 'rnb'
	GenreEDM                          // 'edm'
	GenreBallad                       // 'ballad'
	GenreCountry                      // 'country'
	GenreJazz                         // 'jazz'
	GenreIndie                        // 'indie'
	GenreKPop                         // 'k_pop'
	GenreVPop                         // 'v_pop'
	GenreSoundtrack                   // 'soundtrack'
)
