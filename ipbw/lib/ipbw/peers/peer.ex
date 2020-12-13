defmodule Ipbw.Peers.Peer do
  use Ecto.Schema
  import Ecto.Changeset

  schema "peers" do
    field :ip, :string
    field :pid, :string

    timestamps()
  end

  @doc false
  def changeset(peer, attrs) do
    peer
    |> cast(attrs, [:pid, :ip])
    |> validate_required([:pid, :ip])
  end
end