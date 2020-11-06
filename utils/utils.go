package utils

import (
	"../config"
	"fmt"
)
//Max retourne le max entre 2 horloge logique
func Max(x, y uint32)uint32 {
	if(x>y){
		return x
	}else {
		return y
	}
}

// PrintMessage affiche le message specifiquement formaté pour tous les processus
func PrintMessage(procId uint32, partName string, message string)  {
	if config.DEBUG{
		fmt.Print("[p"+ fmt.Sprint(procId) + partName+ "]" + message)
	}

}