package channelserver

import (
	"fmt"

	"github.com/Andoryuuta/byteframe"
	"github.com/Solenataris/Erupe/network/mhfpacket"
	"go.uber.org/zap"
)

func handleMsgSysEnumerateClient(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgSysEnumerateClient)

	// Read-lock the stages map.
	s.server.stagesLock.RLock()

	stage, ok := s.server.stages[pkt.StageID]
	if !ok {
		char := fmt.Sprintf("(%s: %d)", s.Name, s.charID)
		s.logger.Fatal("Can't enumerate clients for stage that doesn't exist! "+char, zap.String("stageID", pkt.StageID))
	}

	// Unlock the stages map.
	s.server.stagesLock.RUnlock()

	// Read-lock the stage and make the response with all of the charID's in the stage.
	resp := byteframe.NewByteFrame()
	stage.RLock()

	// TODO(Andoryuuta): Is only the reservations needed? Do clients send this packet for mezeporta as well?

	// Make a map to deduplicate the charIDs between the unreserved clients and the reservations.
	deduped := make(map[uint32]interface{})

	// Add the charIDs
	for session := range stage.clients {
		deduped[session.charID] = nil
	}

	for charid := range stage.reservedClientSlots {
		deduped[charid] = nil
	}

	// Write the deduplicated response
	resp.WriteUint16(uint16(len(deduped))) // Client count
	for charid := range deduped {
		resp.WriteUint32(charid)
	}

	stage.RUnlock()

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
	s.logger.Debug("MsgSysEnumerateClient Done!")
}

func handleMsgMhfListMember(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfListMember)

	resp := byteframe.NewByteFrame()
	resp.WriteUint32(0) // Members count. (Unsure of what kind of members these actually are, guild, party, COG subscribers, etc.)

	doAckBufSucceed(s, pkt.AckHandle, resp.Data())
}

func handleMsgMhfOprMember(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfListMember)
	// TODO: add targetid(uint32) to charid(uint32)'s database under new field
	doAckSimpleSucceed(s, pkt.AckHandle, make([]byte, 4))
}

func handleMsgMhfShutClient(s *Session, p mhfpacket.MHFPacket) {}

func handleMsgSysHideClient(s *Session, p mhfpacket.MHFPacket) {}
