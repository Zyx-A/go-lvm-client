package lvm_client

import (
  "errors"
  "strings"
  "strconv"
  "github.com/starkandwayne/go-lvm-client/system"
)

type LogicalVolumeType int
type AllocationPolicyType int

const (  // iota is reset to 0
  LVTUnspecified LogicalVolumeType = iota
  LVTMirrored
  LVTMirroredWithoutInitialSync
  LVTOrigin
  LVTOriginWithMergingSnapshot
  LVTSnapshot
  LVTMergingSnapshot
  LVTPvmove
  LVTVirtual
  LVTMirrorImage
  LVTMirrorImageOutOfSync
  LVTUnderConversion
)

const (
  LVATContiguous AllocationPolicyType = iota
  LVATCling
  LVATNormal
  LVATAnywhere
  LVATInherited
)

type LogicalVolume struct {
  LVName         string
  VGName         string
  Attrs          string
  VolumeType     LogicalVolumeType
  Writable       bool
  AllocationPolicy AllocationPolicyType
  Locked         bool
  LVSize         float64
}

func NewLogicalVolume() LogicalVolume {
  return LogicalVolume{}
}

func (lv *LogicalVolume) ParseLine(lvsLine string, delimiter string) (err error) {
  tokens := strings.Split(strings.Trim(lvsLine, " "), delimiter)
  if (len(tokens) != 10) {
    err = errors.New("Expected 10 colon items from lvs. Perhaps an unsupported operating system.")
    return
  }
  lv.LVName = tokens[0]
  lv.VGName = tokens[1]
  lv.Attrs  = tokens[2]
  lv.parseAttr()

  lv.LVSize, err = strconv.ParseFloat(tokens[3], 32)
  if (err != nil) {
    return err
  }

  return
}

func (lv *LogicalVolume) parseAttr() {
  attrs := strings.Split(lv.Attrs, "")
  switch attrs[0] {
    case "-": lv.VolumeType = LVTUnspecified
    case "m": lv.VolumeType = LVTMirrored
    case "M": lv.VolumeType = LVTMirroredWithoutInitialSync
    case "o": lv.VolumeType = LVTOrigin
    case "O": lv.VolumeType = LVTOriginWithMergingSnapshot
    case "s": lv.VolumeType = LVTSnapshot
    case "S": lv.VolumeType = LVTMergingSnapshot
    case "p": lv.VolumeType = LVTPvmove
    case "v": lv.VolumeType = LVTVirtual
    case "i": lv.VolumeType = LVTMirrorImage
    case "I": lv.VolumeType = LVTMirrorImageOutOfSync
    case "c": lv.VolumeType = LVTUnderConversion
  }
  lv.Writable = attrs[1] == "w"
  switch strings.ToLower(attrs[2]) {
    case "c": lv.AllocationPolicy = LVATContiguous
    case "l": lv.AllocationPolicy = LVATCling
    case "n": lv.AllocationPolicy = LVATNormal
    case "a": lv.AllocationPolicy = LVATAnywhere
    case "i": lv.AllocationPolicy = LVATInherited
  }
  // Capitalised if the volume is currently locked against allocation changes
  lv.Locked = attrs[2] != strings.ToLower(attrs[2])
}

func LogicalVolumes(repo system.SystemRepository) (lvs []LogicalVolume, err error) {
  lvsOutput, delimiter, err := repo.LogicalVolumes()
  lvs = []LogicalVolume{}
  lvsLines := strings.Split(lvsOutput, "\n")
  for _, lvLine := range lvsLines {
    if len(lvLine) > 0 {
      lv := NewLogicalVolume()
      err = lv.ParseLine(lvLine, delimiter)
      if err != nil {
        return
      }
      lvs = append(lvs, lv)
    }
  }

  return
}
